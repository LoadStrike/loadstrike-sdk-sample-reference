from __future__ import annotations

from typing import Any

from playwright.async_api import async_playwright
from selenium.webdriver import ChromeOptions

from loadstrike_sdk import (
    AzureEventHubsEndpointDefinition,
    CrossPlatformScenarioConfigurator,
    DatadogReportingSink,
    DelegateStreamEndpointDefinition,
    GrafanaLokiReportingSink,
    HttpAuthOptions,
    HttpEndpointDefinition,
    HttpOAuth2ClientCredentialsOptions,
    InfluxDbReportingSink,
    KafkaEndpointDefinition,
    KafkaSaslOptions,
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
    NatsEndpointDefinition,
    OtelCollectorReportingSink,
    PushDiffusionEndpointDefinition,
    RabbitMqEndpointDefinition,
    RedisStreamsEndpointDefinition,
    SplunkReportingSink,
    TimescaleDbReportingSink,
    TrackingFieldSelector,
    TrackingPayloadBuilder,
)

FEATURE_OVERVIEW = "This feature area is the smallest end-to-end setup and execution path so readers can confirm package references and register a runner quickly."


def build():
    # Register And Run is configured on a reusable context so you can compare the context flow separately from scenario authoring.
    context = {"runner": _runner(), "context": _runner().build_context()}
    return LoadStrikeRunner.run(context)


class placeholders:
    RUNNER_KEY = "runner_dummy_orders_reference"
    ORDERS_API_BASE_URL = "https://orders.example.test"
    KAFKA_BOOTSTRAP_SERVERS = "localhost:9092"
    NATS_SERVER_URL = "nats://localhost:4222"
    REDIS_CONNECTION_STRING = "localhost:6379,abortConnect=false"
    RABBIT_HOST = "localhost"
    AZURE_EVENT_HUBS_CONNECTION_STRING = (
        "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy"
    )
    REPORT_FOLDER = "./artifacts/reports"
    ORDER_TOPIC = "orders.created"
    EXAMPLE_ORDER_NUMBER = "ORD-10001"
    EXAMPLE_TENANT = "demo-tenant"

    @staticmethod
    def scenario_name(suffix: str) -> str:
        return f"orders.{suffix}"


class OrdersReportingSink:
    """Minimal custom sink showing the public hook names the Python SDK accepts."""

    @property
    def sink_name(self) -> str:
        return "orders-sample-sink"

    SinkName = property(lambda self: self.sink_name)

    def init(self, _context: dict[str, Any], _infra_config: dict[str, Any]) -> None:
        return

    def start(self, _session: dict[str, Any]) -> None:
        return

    def save_realtime_stats(self, _scenario_stats: list[dict[str, Any]]) -> None:
        return

    def save_realtime_metrics(self, _metrics: dict[str, Any]) -> None:
        return

    def save_run_result(self, _result: dict[str, Any]) -> None:
        return

    def stop(self) -> None:
        return


class OrdersRuntimePolicy:
    """Runtime policy example used by the policy-focused feature folders."""

    def should_run_scenario(self, scenario_name: str) -> bool:
        return not scenario_name.endswith(".skip")

    def before_scenario(self, _scenario_name: str) -> None:
        return

    def after_scenario(self, _scenario_name: str, _stats: Any) -> None:
        return

    def before_step(self, _scenario_name: str, _step_name: str) -> None:
        return

    def after_step(self, _scenario_name: str, _step_name: str, _reply: Any) -> None:
        return


class OrdersWorkerPlugin:
    """Worker plugin example showing init, start, get_data, stop, and dispose hooks."""

    plugin_name = "orders-sample-plugin"

    def init(self, _context: Any = None, _infra_config: Any = None) -> None:
        return

    def start(self, _session: Any = None) -> None:
        return

    def get_data(self, _result: Any) -> LoadStrikePluginData:
        table = LoadStrikePluginDataTable.create("Captured Orders")
        table.headers = ["Order Number", "Tenant"]
        table.rows = [
            {
                "Order Number": placeholders.EXAMPLE_ORDER_NUMBER,
                "Tenant": placeholders.EXAMPLE_TENANT,
            }
        ]
        plugin_data = LoadStrikePluginData.create(self.plugin_name)
        plugin_data.hints = ["Reference-only plugin payload for the sample repository."]
        plugin_data.tables = [table]
        return plugin_data

    def stop(self) -> None:
        return

    def dispose(self) -> None:
        return


def _step_reply(context: Any):
    return LoadStrikeStep.run(
        "publish-order",
        context,
        lambda: LoadStrikeResponse.ok(
            {
                "orderNumber": placeholders.EXAMPLE_ORDER_NUMBER,
                "tenant": placeholders.EXAMPLE_TENANT,
            },
            "201",
            128,
            "created",
            4.5,
        ),
    ).as_reply()


def _scenario(name_suffix: str = "publish"):
    return (
        LoadStrikeScenario.create(placeholders.scenario_name(name_suffix), _step_reply)
        .with_init(lambda _context: None)
        .with_clean(lambda _context: None)
        .with_weight(2)
        .without_warm_up()
        .with_max_fail_count(3)
        .with_restart_iteration_on_fail(True)
    )


def _runner():
    return (
        LoadStrikeRunner.register_scenarios(_scenario())
        .with_runner_key(placeholders.RUNNER_KEY)
        .with_test_suite("orders-reference")
        .with_test_name("sample-reference")
        .with_session_id("orders-sample-session")
        .with_report_folder(placeholders.REPORT_FOLDER)
        .with_report_file_name("orders-reference")
        .with_report_formats(LoadStrikeReportFormat.Txt, LoadStrikeReportFormat.Html)
        .with_reporting_interval(2.5)
        .with_logger_config(lambda: {"logger": "dummy"})
        .with_minimum_log_level("Warning")
        .with_runtime_policies(OrdersRuntimePolicy())
        .with_worker_plugins(OrdersWorkerPlugin())
        .with_reporting_sinks(OrdersReportingSink())
        .without_reports()
    )
