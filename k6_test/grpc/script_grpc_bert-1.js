import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load([], 'kfs_inference_v2.proto');
let hostname = __ENV.ENDPOINT===undefined ? 'modelmesh-serving:8033' :  __ENV.ENDPOINT;
let model_name = 'bert-uncased';
let vu = __ENV.K6_VUS===undefined ? 1:__ENV.K6_VUS;
let duration = __ENV.DURATION===undefined ? '1s': __ENV.DURATION;

let params = {
  tags: { model_name: `${model_name}`,},
};

export const options = {
  discardResponseBodies: true,
  scenarios: {
    constant_vus: {
      executor: 'constant-vus',
      vus: `${vu}`,
      duration: `${duration}`,
    },
//    scen2: {
//      executor: 'shared-iterations',
//      vus: 10,
//      iterations: 200
//    },
  },
};
const data = {"model_name":`${model_name}`, "inputs":[{"name":"input_word_ids", "shape": [1, 16], "datatype":"INT32","contents":{"int_contents":[101,1045,2293,13137,20968,102,1045,2293,13137,20968,102,0,0,0,0,0]}}, {"name":"input_mask", "shape":[1,16], "datatype":"INT32","contents":{"int_contents":[1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0]}} ,{"name": "input_type_ids","shape":[1,16], "datatype":"INT32", "contents":{"int_contents":[0,0,0,0,0,0,1,1,1,1,1,0,0,0,0,0]}}]};

export default function (){
  client.connect(`${hostname}`, { plaintext: true });
  const res= client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);
  check(res, {
    'status is OK': (res) => res && res.status == grpc.StatusOK,
  });
//  console.log(JSON.stringify(res.message));
  client.close();
};
