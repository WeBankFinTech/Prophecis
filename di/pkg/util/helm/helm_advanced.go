// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"webank/DI/commons/service"

	log "github.com/sirupsen/logrus"
)

/*
* Generate value file
 */
func GenerateValueFile(values interface{}) (valueFileName string, err error) {
	// 1. generate the template file
	valueFile, err := ioutil.TempFile(os.TempDir(), "values")
	if err != nil {
		log.Errorf("Failed to create tmp file %v due to %v", valueFile.Name(), err)
		return "", err
	}

	valueFileName = valueFile.Name()
	log.Debugf("Save the values file %s", valueFileName)

	// 2. dump the object into the template file
	err = toYaml(values, valueFile)
	return valueFileName, err
}

/**
* generate helm template without tiller: helm template -f values.yaml chart_name
* Exec /usr/local/bin/helm, [template -f /tmp/values313606961 --namespace default --name hj /charts/tf-horovod]
* returns generated template file: templateFileName
 */
func GenerateHelmTemplate(name string, namespace string, valueFileName string, chartName string, command string, req *service.JobDeploymentRequest) (templateFileName string, err error) {
	tempName := fmt.Sprintf("%s.yaml", name)
	templateFile, err := ioutil.TempFile("", tempName)
	if err != nil {
		return templateFileName, err
	}
	templateFileName = templateFile.Name()

	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return templateFileName, err
	}

	// 3. check if the chart file exists
	// if strings.HasPrefix(chartName, "/") {
	if _, err = os.Stat(chartName); err != nil {
		return templateFileName, err
	}
	// }

	// 4. prepare the arguments
	args := []string{binary, "template", "-f", valueFileName,
		"--namespace", namespace,
		"--name", name, chartName, ">", templateFileName}

	log.Debugf("Exec bash -c %v", args)

	// return syscall.Exec(cmd, args, env)
	// 5. execute the command
	log.Debugf("Generating template  %v", args)
	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	// cmd.Env = env
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", string(out))
	if err != nil {
		return templateFileName, fmt.Errorf("failed to execute %s, %v with %v", binary, args, err)
	}

	// // 6. clean up the value file if needed
	// if log.GetLevel() != log.DebugLevel {
	// 	err = os.Remove(valueFileName)
	// 	if err != nil {
	// 		log.Warnf("Failed to delete %s due to %v", valueFileName, err)
	// 	}
	// }

	//FIXME FIX cmd for TFJob.yaml
	err = RecreateTFJobYAML(templateFileName, command, name, req)
	if err != nil {
		return templateFileName, fmt.Errorf("failed to ReplaceCMD %v", command)
	}

	return templateFileName, nil
}

func RecreateTFJobYAML(path string, cmd string, name string, req *service.JobDeploymentRequest) error {

	replace, _ := myEscape(cmd)
	replace = "                - \"" + replace + "\""

	//get cmd line number
	cmdLineNum, _ := getLine(path, "              command:")
	envLineNum, _ := getLine(path, "              env:")

	startline := cmdLineNum + 3
	count := envLineNum - startline
	log.Printf("cmdLineNum: %v, envLineNum: %v,startline: %v,count: %v", cmdLineNum, envLineNum, startline, count)

	//assemble cmd
	AssembleCMD(path, replace, name, req)

	return nil
}

func myEscape(source string) (string, error) {
	//var j int = 0
	if len(source) == 0 {
		return "", errors.New("source is null")
	}
	tempStr := source[:]
	desc := make([]string, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		//flag := false
		var escape string
		switch tempStr[i] {
		case '\r':
			//flag = true
			escape = `\r`
			break
		case '\n':
			//flag = true
			escape = `\n`
			break
		case '\t':
			//flag = true
			escape = `\t`
			break
		//case ':':
		//        //flag = true
		//        escape = `\:`
		//        break
		case '"':
			//flag = true
			escape = `\"`
			break
		case '\\':
			//flag = true
			escape = `\\`
			break
		default:
			escape = string(tempStr[i])
		}
		//if flag {
		//        desc[j] = '\\'
		//        desc[j+1] = escape
		//        j = j + 2
		//} else {
		//        desc[j] = tempStr[i]
		//        j = j + 1
		//}
		desc = append(desc, escape)
	}
	return strings.Join(desc, ""), nil
}

func getLine(path string, str string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			return line, nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}
	return 0, nil
}

func AssembleCMD(path string, cmd string, name string, req *service.JobDeploymentRequest) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	//添加subPath行
	//获取,                  - name: "static-volume-1"
	//var vmstaticline []int
	var psCPURequests = 0
	var psCPULimits = 0
	cpu := req.PSCPU
	for i, line := range lines {
		if strings.Contains(line, `                  mountPath: "/job"`) {
			lines[i] = line + "\n                  subPath: " + name
		}

		if strings.Contains(line, "              command:") {
			lines[i+1] = "                - bash"
		}

		if strings.Contains(line, `              requests:`) && psCPURequests == 0 {
			lines[i] = line + "\n                  cpu: " + cpu
			psCPURequests += 1
		}
		if strings.Contains(line, `              limits:`) && psCPULimits == 0 {
			lines[i] = line + "\n                  cpu: " + cpu
			psCPULimits += 1
		}
	}

	newLines := ReplaceCMD(lines, cmd)

	//log.Debugf("AssembleCMD, print old lines")
	//for _, line := range lines {
	//	log.Println(line)
	//}

	//log.Debugf("AssembleCMD, print new lines")
	//for _, line := range newLines {
	//	log.Println(line)
	//}

	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func ReplaceCMD(lines []string, cmd string) []string {

	var envlines, startlines []int

	for i, line := range lines {
		//fmt.Println(line)

		if strings.Contains(line, `              command:`) {
			startlines = append(startlines, i)
		}
		if strings.Contains(line, `              env:`) {
			envlines = append(envlines, i)
		}
	}

	fmt.Println(startlines, envlines)

	var newLines []string

	for i, v := range startlines {
		startLineNum := v + 3
		envLineNum := envlines[i]

		if i == 0 {
			newLines = append(lines[:0], lines[0:startLineNum]...)
		} else {
			lastEnvLineNum := envlines[i-1]
			newLines = append(newLines, lines[lastEnvLineNum:startLineNum]...)
		}

		newLines = append(newLines, cmd)

		if i == len(envlines)-1 {
			newLines = append(newLines, lines[envLineNum:]...)
		}

	}

	//log.Debugf("AssembleCMD, print old lines")
	//for _, line := range lines {
	//	log.Println(line)
	//}

	return newLines

}

/**
* Check the chart version by given the chart directory
* helm inspect chart /charts/tf-horovod
 */

func GetChartVersion(chart string) (version string, err error) {
	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return "", err
	}

	// 1. check if the chart file exists, if it's it's unix path, then check if it's exist
	// if strings.HasPrefix(chart, "/") {
	if _, err = os.Stat(chart); err != nil {
		return "", err
	}
	// }

	// 2. prepare the arguments
	args := []string{binary, "inspect", "chart", chart,
		"|", "grep", "version:"}
	log.Debugf("Exec bash -c %v", args)

	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("Failed to find version when executing %s, result is %s", args, out)
	}
	fields := strings.Split(lines[0], ":")
	if len(fields) != 2 {
		return "", fmt.Errorf("Failed to find version when executing %s, result is %s", args, out)
	}

	version = strings.TrimSpace(fields[1])
	return version, nil
}

func GetChartName(chart string) string {
	return filepath.Base(chart)
}
