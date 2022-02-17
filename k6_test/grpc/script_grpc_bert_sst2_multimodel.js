import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load([], 'kfs_inference_v2.proto');
const inputs = [
{"name":"input_word_ids",
"shape": [1, 128],
"datatype": "INT32",
"contents": {"int_contents":[  101,  2182,  2057,  2175,  2059,  1010,  2017,  1998,  1045,  2003,  1037,  2639,  2201,  2011,  5046,  3769,  3063, 22294,  2368, 16768,  1012,  2009,  2001, 16768,  1005,  1055,  2117,  3729,  2004,  1037,  3948,  3063,  1012,   102, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,0,0,0]}
},
{"name":"input_mask",
"shape":[1,128],
"datatype":"INT32",
"contents":{"int_contents":[1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]}
},
{"name": "input_type_ids",
"shape":[1,128],
"datatype":"INT32",
"contents":{"int_contents":[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]}
},
]

export default () => {
  client.connect('modelmesh-serving:8033', { plaintext: true });
  //send request to multiple models at once
  for (let i = 1; i <= 2; i++) {
    let params = {
      tags: { model_name: `experts-bert-wiki-books-sst2-${i}` },
    };
    const data = {"model_name":`experts-bert-wiki-books-sst2-${i}`, "inputs":inputs}
    let r = client.invoke('inference.GRPCInferenceService/ModelInfer', data, params);
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    //console.log(JSON.stringify(r.message));
  }
  client.close();
};
