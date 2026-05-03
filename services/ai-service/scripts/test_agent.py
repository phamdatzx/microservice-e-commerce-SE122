"""
Interactive test for the e-commerce AI agent.

Runs a few predefined queries and then drops into an interactive loop
so you can chat freely.

Usage:
    python scripts/test_agent.py
"""

from __future__ import annotations

import sys
from pathlib import Path

# Ensure app.* imports work from scripts/
sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)

from app.agent.agent import chat


def divider(label: str = "") -> None:
    print(f"\n{'═' * 60}")
    if label:
        print(f"  {label}")
        print(f"{'═' * 60}")


def run_test(query: str) -> None:
    divider(f"Query: {query}")
    response = chat(query)
    print(f"\n🤖 Response:\n{response}\n")


def main() -> None:
    # ── Predefined test queries ───────────────────────────────────────
    test_queries = [
        "Tôi muốn mua máy giặt, có sản phẩm nào không?",
        "Cho tôi xem đơn hàng đang giao",
        "Tôi có voucher nào không?",
    ]

    print("🚀 E-commerce AI Agent — Test Suite")
    divider("Running predefined queries")

    for q in test_queries:
        run_test(q)

    # ── Interactive mode ──────────────────────────────────────────────
    divider("Interactive mode (type 'quit' to exit)")

    while True:
        try:
            user_input = input("\n👤 You: ").strip()
        except (EOFError, KeyboardInterrupt):
            print("\n👋 Bye!")
            break

        if not user_input:
            continue
        if user_input.lower() in {"quit", "exit", "q"}:
            print("👋 Bye!")
            break

        response = chat(user_input)
        print(f"\n🤖 Agent: {response}")


if __name__ == "__main__":
    main()
