export JAVA_HOME=
source env.sh #HDFS ENV Config

# Tensorflow HDFS ENV
source ${HADOOP_HOME}/libexec/hadoop-config.sh
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${JAVA_HOME}/jre/lib/amd64/server
export CLASSPATH=$(${HADOOP_HOME}/bin/hadoop classpath --glob)


# Pyspark ENV
export PYSPARK_PYTHON=${ANDCONDA_HOME}/bin/python
export PYSPARK_DRIVER_PYTHON=/opt/conda/bin/python