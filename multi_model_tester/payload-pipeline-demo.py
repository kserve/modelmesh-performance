from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_payload_op = components.load_component_from_file("payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-paylood-automation",
    description="Model mesh payload automation",
)
def modelmesh_pipeline(
    payload_model_type: str = "SimpleStringTF CifarPytorch MnistOnnx",
    payload_query_duration_in_seconds: str = "60"
):
    task_payload = task_payload_op(model_type=payload_model_type,
                                   inference_server_url="dns:///modelmesh-serving.modelmesh-serving:8033",
                                   number_of_models="5000",
                                   number_of_workers="100",
                                   duration=payload_query_duration_in_seconds,
                                   queries_per_second="400",
                                   number_of_connections="100")
    task_payload.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh-payload" + ".yaml")
