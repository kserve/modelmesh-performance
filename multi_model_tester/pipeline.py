from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_deploy_op = components.load_component_from_file("deploy-component.yaml")
task_payload_op = components.load_component_from_file("payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-automation",
    description="Model mesh deployment and payload automation",
)
def modelmesh_pipeline(
):
    task_deploy = task_deploy_op()
    task_payload = task_payload_op().after(task_deploy)
    task_deploy.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')
    task_payload.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh" + ".yaml")
