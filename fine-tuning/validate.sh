#!/bin/bash
docker model package --gguf $PWD/output/fine-tuned-model.F16.gguf demo/queen-pedauque:0.5b-0.0.0
docker model list
docker model run demo/queen-pedauque:0.5b-0.0.0
