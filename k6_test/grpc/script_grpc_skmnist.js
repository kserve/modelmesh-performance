import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

{{k6_opts}}

const client = new grpc.Client();
client.load([], '../k6_test/kfs_inference_v2.proto');
const inputsData = JSON.parse(open(`../k6_test/payloads/{{payload}}`));
let params = {
  tags: { model_name: `snapml-mnist-svm` },
}
export default () => {
  client.connect('{{base_url}}', { plaintext: true });
  const data = { 
    "model_name": "snapml-mnist-svm",
    "inputs": inputsData["inputs"]
  };
  const response = client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);

  check(response, {
    'status is OK': (response) => response && response.status === grpc.StatusOK,
  });
 // console.log(JSON.stringify(response.message));

  client.close();
  sleep(1);
};
