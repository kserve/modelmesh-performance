import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
let model_name = 'bert-uncased'
client.load([], 'kfs_inference_v2.proto');
    let params = {
      tags: { model_name: `${model_name}` },
    };
export default () => {
  client.connect('modelmesh-serving:8033', { plaintext: true });
  const data = {"model_name":`${model_name}`, "inputs":[{"name":"input_word_ids", "shape": [1, 16], "datatype":"INT32","contents":{"int_contents":[101,1045,2293,13137,20968,102,1045,2293,13137,20968,102,0,0,0,0,0]}}, {"name":"input_mask", "shape":[1,16], "datatype":"INT32","contents":{"int_contents":[1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0]}} ,{"name": "input_type_ids","shape":[1,16], "datatype":"INT32", "contents":{"int_contents":[0,0,0,0,0,0,1,1,1,1,1,0,0,0,0,0]}}]}
  const res= client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);

  check(res, {
    'status is OK': (res) => res && res.status == grpc.StatusOK,
  });
//  console.log(JSON.stringify(res.message));

  client.close();
//  sleep(1);
};
