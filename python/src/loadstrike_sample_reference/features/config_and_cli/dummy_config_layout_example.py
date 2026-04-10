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

FEATURE_OVERVIEW = "This feature area focuses on loading config, overriding config with run arguments, and documenting the expected config layout."


def build():
    # Dummy Config Layout keeps the file-based setup visible while the runtime path stays minimal.
    config_layout = {"ConfigFile": "loadstrike.config.json", "InfraFile": "loadstrike.infra.json", "RunnerKey": placeholders.RUNNER_KEY}
    context = (LoadStrikeRunner.register_scenarios(_scenario("config-layout")).load_config(f"features/config_and_cli/{config_layout['ConfigFile']}").load_infra_config(f"features/config_and_cli/{config_layout['InfraFile']}").with_runner_key(config_layout["RunnerKey"]).build_context())
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
