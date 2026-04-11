import { Agent, fetch } from "undici";
import { LoadStrikeResponse } from "@loadstrike/loadstrike-sdk";

export const sharedHttpClient = new Agent({
  connections: 1000,
  keepAliveTimeout: 120_000,
  keepAliveMaxTimeout: 300_000,
  connectTimeout: 15_000,
  pipelining: 1,
});

export async function postOrder(orderId, tenantId) {
  const response = await fetch("https://api.example.com/orders", {
    method: "POST",
    headers: {
      "content-type": "application/json",
      "accept-encoding": "gzip, deflate",
    },
    body: JSON.stringify({
      orderId,
      tenantId,
      amount: 49.95,
    }),
    dispatcher: sharedHttpClient,
  });

  await response.arrayBuffer();

  return response.ok
    ? LoadStrikeResponse.ok(String(response.status))
    : LoadStrikeResponse.fail(String(response.status));
}

export async function closeSharedHttpClient() {
  await sharedHttpClient.close();
}
