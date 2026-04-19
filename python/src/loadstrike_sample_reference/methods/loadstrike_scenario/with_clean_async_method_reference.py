from __future__ import annotations

from typing import Any

from loadstrike_sdk import (
    CorrelationStoreConfiguration,
    CrossPlatformScenarioConfigurator,
    HttpEndpointDefinition,
    LoadStrikeMetric,
    LoadStrikePluginData,
    LoadStrikePluginDataTable,
    LoadStrikeReportFormat,
    LoadStrikeResponse,
    LoadStrikeRunner,
    LoadStrikeScenario,
    LoadStrikeSimulation,
    LoadStrikeStep,
    LoadStrikeThreshold,
    RedisCorrelationStoreOptions,
    TrackingFieldSelector,
    TrackingPayloadBuilder,
)

RUNNER_KEY = "runner_dummy_orders_reference"

class TempConfigPaths:
    def __init__(self, config_path: str, infra_path: str) -> None:
        self.config_path = config_path
        self.infra_path = infra_path


class OrdersReportingSink:
    @property
    def sink_name(self) -> str:
        return "orders-sample-sink"

    def init(self, _context: Any, _infra_config: Any) -> None:
        return

    def start(self, _session: Any) -> None:
        return

    def save_realtime_stats(self, _stats: Any) -> None:
        return

    def save_realtime_metrics(self, _metrics: Any) -> None:
        return

    def save_run_result(self, _result: Any) -> None:
        return

    def stop(self) -> None:
        return


class OrdersRuntimePolicy:
    def should_run_scenario(self, _scenario_name: str) -> bool:
        return True

    def before_scenario(self, _scenario_name: str) -> None:
        return

    def after_scenario(self, _scenario_name: str, _stats: Any) -> None:
        return

    def before_step(self, _scenario_name: str, _step_name: str) -> None:
        return

    def after_step(self, _scenario_name: str, _step_name: str, _reply: Any) -> None:
        return


class OrdersWorkerPlugin:
    plugin_name = "orders-sample-plugin"

    def init(self, _context: Any = None, _infra_config: Any = None) -> None:
        return

    def start(self, _session: Any = None) -> None:
        return

    def get_data(self, _result: Any) -> LoadStrikePluginData:
        return LoadStrikePluginData.create(self.plugin_name)

    def stop(self) -> None:
        return

    def dispose(self) -> None:
        return


def create_order_reply():
    return LoadStrikeResponse.ok("200", 128, "ok", 3.2)


def execute_order_get(context: Any):
    return LoadStrikeStep.run("get-order", context, lambda: create_order_reply()).as_reply()


def baseline_scenario(name: str = "orders.get-by-id"):
    return (
        LoadStrikeScenario.create(name, execute_order_get)
        .with_load_simulations(LoadStrikeSimulation.iterations_for_constant(1, 1))
        .without_warm_up()
    )


def base_runner():
    return (
        LoadStrikeRunner.create()
        .add_scenario(baseline_scenario())
        .with_runner_key(RUNNER_KEY)
        .with_test_suite("orders-reference")
        .with_test_name("orders-get-by-id")
        .without_reports()
    )


def base_context():
    return base_runner().build_context()


def http_source():
    return HttpEndpointDefinition(
        name="orders-http-source",
        mode="Produce",
        tracking_field=TrackingFieldSelector.parse("header:X-Correlation-Id"),
        url="https://orders.example.test/api/orders",
        method="GET",
        response_source="Body",
    )


def http_destination():
    return HttpEndpointDefinition(
        name="orders-http-destination",
        mode="Consume",
        tracking_field=TrackingFieldSelector.parse("json:$.trackingId"),
        gather_by_field=TrackingFieldSelector.parse("json:$.tenantId"),
        url="https://orders.example.test/api/order-events",
        method="GET",
        response_source="Body",
        consume_json_array_response=True,
        consume_array_path="$.items",
    )


def tracking_configuration():
    return {
        "Source": http_source(),
        "Destination": http_destination(),
        "RunMode": "GenerateAndCorrelate",
        "CorrelationTimeoutSeconds": 30,
        "TimeoutSweepIntervalSeconds": 1,
        "TimeoutBatchSize": 200,
        "TimeoutCountsAsFailure": True,
        "MetricPrefix": "orders_tracking",
        "ExecuteOriginalScenarioRun": False,
        "CorrelationStore": CorrelationStoreConfiguration.in_memory(),
    }


def tracked_context():
    return (
        LoadStrikeRunner.create()
        .add_scenario(
            CrossPlatformScenarioConfigurator.configure(
                baseline_scenario("orders.tracked"),
                tracking_configuration(),
            )
        )
        .with_runner_key(RUNNER_KEY)
        .without_reports()
        .build_context()
    )


def build_tracking_payload():
    builder = TrackingPayloadBuilder()
    builder.set_body('{"trackingId":"ord-1"}')
    return builder.build()


def run_result():
    return None


def scenario_stats():
    return None


def write_temp_config_files():
    return TempConfigPaths(
        "method-reference.loadstrike.config.json",
        "method-reference.loadstrike.infra.json",
    )

class WithCleanAsyncMethodReference:
    @staticmethod
    def attach_async_no_op_hook_example() -> Any:
        """Attach an asynchronous lifecycle hook with no extra side effects."""
        return baseline_scenario().with_clean_async(lambda context: None)

    @staticmethod
    def attach_async_metric_hook_example() -> Any:
        """Attach an asynchronous lifecycle hook that still registers a metric."""
        return baseline_scenario().with_clean_async(lambda context: context.register_metric(LoadStrikeMetric.counter("orders_seen", "count")))
