import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load([], 'kfs_inference_v2.proto');
let params = {
  tags: { model_name: `example-mnist-predictor` },
}
export default () => {
  client.connect('modelmesh-serving:8033', { plaintext: true });
  const data = { "model_name": "example-mnist-predictor", "inputs": [{ "name": "predict", "shape": [1, 64], "datatype": "FP32", "contents": { "fp32_contents": [0.0, 0.0, 1.0, 11.0, 14.0, 15.0, 3.0, 0.0, 0.0, 1.0, 13.0, 16.0, 12.0, 16.0, 8.0, 0.0, 0.0, 8.0, 16.0, 4.0, 6.0, 16.0, 5.0, 0.0, 0.0, 5.0, 15.0, 11.0, 13.0, 14.0, 0.0, 0.0, 0.0, 0.0, 2.0, 12.0, 16.0, 13.0, 0.0, 0.0, 0.0, 0.0, 0.0, 13.0, 16.0, 16.0, 6.0, 0.0, 0.0, 0.0, 0.0, 16.0, 16.0, 16.0, 7.0, 0.0, 0.0, 0.0, 0.0, 11.0, 13.0, 12.0, 1.0, 0.0] }}]};
  const response = client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);

  check(response, {
    'status is OK': (response) => response && response.status === grpc.StatusOK,
  });
 // console.log(JSON.stringify(response.message));

  client.close();
  sleep(1);
};
