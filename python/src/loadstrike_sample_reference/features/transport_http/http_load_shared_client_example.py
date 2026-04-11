from __future__ import annotations

from loadstrike_sdk import (
    LoadStrikeRunner,
    LoadStrikeScenario,
    LoadStrikeSimulation,
    LoadStrikeStep,
)

from .http_load_shared_client_support import close_shared_session, post_order


def build():
    scenario = (
        LoadStrikeScenario.create("http-load-shared-client", _submit_orders)
        .with_load_simulations(LoadStrikeSimulation.inject(25, 1, 30))
    )

    try:
        result = (
            LoadStrikeRunner.register_scenarios(scenario)
            .with_runner_key("runner_dummy_orders_reference")
            .with_test_suite("orders-reference")
            .with_test_name("http-load-shared-client")
            .without_reports()
            .run()
        )
        return {
            "note": "Close the shared HTTP session when the run ends.",
            "Result": result,
        }
    finally:
        close_shared_session()


def _submit_orders(context):
    context.scenario_instance_data.setdefault("tenantId", "tenant-a")
    order_id = f"ord-{context.invocation_number}"

    return LoadStrikeStep.run(
        "POST /orders",
        context,
        lambda: post_order(order_id, context.scenario_instance_data["tenantId"]),
    ).as_reply()
