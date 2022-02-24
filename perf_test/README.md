# Performance Test

WIP

From the top directory:
- Install Python dependencies with

```python
pip3 install -r perf_test/requirements.txt
```

- Run Howitzer:

```sh
export TEMPLATE_DIR=k6_test/
export CONFIG_FILE=perf_test/configs/inference-test.json
./perf_test/runHowitzer.sh
```