from __future__ import annotations

from requests import Session
from requests.adapters import HTTPAdapter

from loadstrike_sdk import LoadStrikeResponse

_shared_session = Session()
_shared_adapter = HTTPAdapter(pool_connections=128, pool_maxsize=1000, max_retries=0)
_shared_session.mount("https://", _shared_adapter)
_shared_session.mount("http://", _shared_adapter)
_shared_session.headers.update({"Accept-Encoding": "gzip, deflate"})


def post_order(order_id: str, tenant_id: str):
    response = _shared_session.post(
        "https://api.example.com/orders",
        json={
            "orderId": order_id,
            "tenantId": tenant_id,
            "amount": 49.95,
        },
        timeout=15,
    )

    if response.ok:
        return LoadStrikeResponse.ok(str(response.status_code))
    return LoadStrikeResponse.fail(str(response.status_code))


def close_shared_session() -> None:
    _shared_session.close()
