from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_payload_op = components.load_component_from_file("payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-paylood-automation",
    description="Model mesh payload automation",
)
def modelmesh_pipeline(
    payload_model_type: str = "MnistSklearn",
    payload_inference_server_url: str = "dns:///modelmesh-serving.modelmesh-serving:8033",
    payload_model_num: str = "1",
    payload_worker_num: str = "1",
    payload_query_duration_in_second: str = "1",
    payload_queries_per_second: str = "1",
    payload_max_connection_num: str = "1"
):
    task_payload = task_payload_op(model_type=payload_model_type,
                                   inference_server_url=payload_inference_server_url,
                                   number_of_models=payload_model_num,
                                   number_of_workers=payload_worker_num,
                                   duration=payload_query_duration_in_second,
                                   queries_per_second=payload_queries_per_second,
                                   number_of_connections=payload_max_connection_num)
    task_payload.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh-payload" + ".yaml")
