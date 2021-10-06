from kfp import dsl, components
from kfp_tekton.compiler import TektonCompiler

task_delete_op = components.load_component_from_file("delete-component.yaml")
task_payload_op = components.load_component_from_file("payload-component.yaml")


@dsl.pipeline(
    name="modelmesh-deleting",
    description="Model mesh delete automation",
)
def modelmesh_pipeline(
    delete_prefix: str = "mnist-sklearn"
):
    task_delete = task_delete_op(prefix=delete_prefix)
    task_delete.add_pod_label('pipelines.kubeflow.org/cache_enabled', 'false')


if __name__ == "__main__":
    TektonCompiler().compile(modelmesh_pipeline, "modelmesh-delete" + ".yaml")
