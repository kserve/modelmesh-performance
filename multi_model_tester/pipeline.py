from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_deploy_op = components.load_component_from_file("deploy-component.yaml")
task_payload_op = components.load_component_from_file("payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-automation",
    description="Model mesh deployment and payload automation",
)
def modelmesh_pipeline(
    num_of_additional_predictors: str = "10",
    deployment_model_type: str = "SimpleStringTF",
    payload_model_type: str = "SimpleStringTF CifarPytorch MnistOnnx",
    payload_query_duration_in_seconds: str = "60"
):
    task_deploy = task_deploy_op(num_of_additional_predictors=num_of_additional_predictors,
                                 model_type=deployment_model_type)
    task_deploy.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')
    task_payload = task_payload_op(model_type=payload_model_type,
                                   inference_server_url="dns:///modelmesh-serving.modelmesh-serving:8033",
                                   number_of_models="5000",
                                   number_of_workers="100",
                                   duration=payload_query_duration_in_seconds,
                                   queries_per_second="400",
                                   number_of_connections="100")
    task_payload.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false').after(task_deploy)


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh" + ".yaml")
