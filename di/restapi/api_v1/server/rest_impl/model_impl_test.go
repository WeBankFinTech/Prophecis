package rest_impl

import "testing"

func TestCommandTransform(t *testing.T) {
	CommandTransform("python3 ${MODEL_DIR}/tf-model/convolutional_network.py --run_date ${run_date} --trainImagesFile ${DATA_DIR}/train-images-idx3-ubyte.gz --trainLabelsFile ${DATA_DIR}/train-labels-idx1-ubyte.gz --testImagesFile ${DATA_DIR}/t10k-images-idx3-ubyte.gz --testLabelsFile ${DATA_DIR}/t10k-labels-idx1-ubyte.gz --learningRate 0.001 --trainingIters 2000\n  ")
	//print("123")

}