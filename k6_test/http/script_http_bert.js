import http from 'k6/http';
import { check } from 'k6';

const data = {"inputs":[{"name":"input_word_ids", "shape": [1, 16], "datatype":"INT32", "data":[101,1045,2293,13137,20968,102,1045,2293,13137,20968,102,0,0,0,0,0]}, {"name":"input_mask", "shape":[1,16], "datatype":"INT32", "data":[1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0]} ,{"name": "input_type_ids","shape":[1,16], "datatype":"INT32", "data":[0,0,0,0,0,0,1,1,1,1,1,0,0,0,0,0]}]}

const model_name = 'bert-uncased';
let params = {headers: { 'Content-Type': 'application/json' },tags: { model_name: `${model_name}` }};

export default () =>{
  const url = `http://modelmesh-serving:8008/v2/models/${model_name}/infer`;
  let res = http.post(url, JSON.stringify(data), params);
  check(res, { 'status is 200': (res) => res.status == 200});
 // console.log(JSON.stringify(res.message))
};
