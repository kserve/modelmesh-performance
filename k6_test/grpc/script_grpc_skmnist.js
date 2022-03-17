import grpc from 'k6/net/grpc';
import { check } from 'k6';
import execution from 'k6/execution';

{{k6_opts}}

const client = new grpc.Client();
client.load([], '../k6_test/kfs_inference_v2.proto');
const inputsData = JSON.parse(open(`../k6_test/payloads/{{payload}}`));
let params = {
  tags: { model_name: `{{model_name}}` },
}

export function setup(){
// Abort on connection errors
  try {
    client.connect('{{base_url}}', { plaintext: true});
  } catch (error) {
    check(error, {"Setup error": (error) => error === null})
    execution.test.abort(error);
  }
}

export default () => {
  client.connect('{{base_url}}', { plaintext: true });
  const data = {
    "model_name": "{{model_name}}",
    "inputs": inputsData["inputs"]
  };
  const response = client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);

  check(response, {
    'status is OK': (response) => response && response.status === grpc.StatusOK,
  });

  client.close();
};
