import { check } from "k6";
import { Counter } from "k6/metrics";
import grpc from "k6/net/grpc";

const grpcReqs = new Counter("grpc_reqs");

function getClient(protoFilePath) {
  const client = new grpc.Client();

  client.load([], protoFilePath);

  return client;
}

function checkResponse(res) {
  check(res, {
    "status is OK": (r) => r && r.status === grpc.StatusOK,
  });
}

export class GrpcClient {
  constructor(options) {
    this.grpcHost = options.grpcHost;
    this.client = getClient(options.protoFilePath);
    this.inferRPCName = options.inferRPCName;

    // Client can't connect on the init context
    this.connected = false;
  }

  infer(data, params) {
    if (!this.connected) {
      this.client.connect(this.grpcHost, { plaintext: true });
      this.connected = true;
    }

    const res = this.client.invoke(this.inferRPCName, data, params);
    checkResponse(res);
    grpcReqs.add(1);
  }

  close() {
    this.client.close();
  }
}