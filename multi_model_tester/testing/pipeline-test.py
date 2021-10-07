from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_deploy_op = components.load_component_from_file("deploy-test-component.yaml")
task_payload_op = components.load_component_from_file("../payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-automation",
    description="Model mesh deployment and payload automation",
)
def modelmesh_pipeline(
    deployment_worker_num: str = "1",
    deployment_prefix: str = "mnist-sklearn",
    deployment_startID: str = "1",
    deployment_endID: str = "1",
    deployment_script_name: str = "deploy_1mnist_sklearn_predictor.sh",
    payload_model_type: str = "MnistSklearn",
    payload_inference_server_url: str = "dns:///modelmesh-serving.modelmesh-serving:8033",
    payload_model_num: str = "1",
    payload_worker_num: str = "1",
    payload_query_duration_in_second: str = "1",
    payload_queries_per_second: str = "1",
    payload_max_connection_num: str = "1"
):
    task_deploy = task_deploy_op(parallel=deployment_worker_num,
                                 prefix=deployment_prefix,
                                 startid=deployment_startID,
                                 endid=deployment_endID,
                                 script=deployment_script_name)
    task_payload = task_payload_op(model_type=payload_model_type,
                                   inference_server_url=payload_inference_server_url,
                                   number_of_models=payload_model_num,
                                   number_of_workers=payload_worker_num,
                                   duration=payload_query_duration_in_second,
                                   queries_per_second=payload_queries_per_second,
                                   number_of_connections=payload_max_connection_num).after(task_deploy)
    task_deploy.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')
    task_payload.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh" + ".yaml")
