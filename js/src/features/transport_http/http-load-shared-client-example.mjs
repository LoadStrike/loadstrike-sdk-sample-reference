import {
  LoadStrikeRunner,
  LoadStrikeScenario,
  LoadStrikeSimulation,
  LoadStrikeStep,
} from "@loadstrike/loadstrike-sdk";
import {
  closeSharedHttpClient,
  postOrder,
} from "./http-load-shared-client-support.mjs";

export async function build() {
  const scenario = LoadStrikeScenario
    .create("http-load-shared-client", async (context) => {
      context.scenarioInstanceData.tenantId ??= "tenant-a";
      const orderId = `ord-${context.invocationNumber}`;

      const step = await LoadStrikeStep.run("POST /orders", context, async () =>
        postOrder(orderId, String(context.scenarioInstanceData.tenantId)));

      return step.asReply();
    })
    .withLoadSimulations(LoadStrikeSimulation.inject(25, 1, 30));

  try {
    const result = await LoadStrikeRunner
      .registerScenarios(scenario)
      .withRunnerKey("runner_dummy_orders_reference")
      .withTestSuite("orders-reference")
      .withTestName("http-load-shared-client")
      .withoutReports()
      .run();

    return {
      note: "Close the shared HTTP client when the run ends.",
      result,
    };
  } finally {
    await closeSharedHttpClient();
  }
}
