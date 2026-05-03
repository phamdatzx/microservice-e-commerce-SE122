"""
test_pipeline_local.py
======================
Test kết hợp Intent Classifier + NER model.

Input  : Câu hỏi của khách hàng
Output : Route name (intent) + Extracted keywords (NER entities)
         → Mô phỏng đầu vào cho Router của chatbot service.

Cách chạy:
    cd services/ai-service
    python test_pipeline_local.py
"""

import os
import sys
import logging
import warnings
import torch
import pandas as pd
from transformers import (
    AutoTokenizer,
    AutoModelForSequenceClassification,
    pipeline as hf_pipeline,
)
from sklearn.preprocessing import LabelEncoder
from underthesea import word_tokenize

warnings.filterwarnings("ignore")

# Đảm bảo có thể import từ thư mục gốc của ai-service
from pathlib import Path
BASE_DIR = str(Path(__file__).resolve().parents[2])
if BASE_DIR not in sys.path:
    sys.path.insert(0, BASE_DIR)

from app.core.config import get_settings

sys.stdout.reconfigure(encoding="utf-8")
logging.basicConfig(level=logging.WARNING, format="%(levelname)s: %(message)s")

settings = get_settings()

INTENT_ROUTE_MAP = {
    "search_product":             "SEARCH_PRODUCT",
    "product_detail":             "PRODUCT_DETAIL",
    "review_summary":             "REVIEW_SUMMARY",
    "compare_product":            "COMPARE_PRODUCT",
    "check_order_status":         "CHECK_ORDER_STATUS",
    "product_voucher":            "PRODUCT_VOUCHER",
    "seller_voucher":             "SELLER_VOUCHER",
    "product_applicable_voucher": "PRODUCT_APPLICABLE_VOUCHER",
    "seller_category":            "SELLER_CATEGORY",
    "greeting":                   "GREETING",
    "goodbye":                    "GOODBYE",
    "out_of_scope":               "OUT_OF_SCOPE",
}

# Màu ANSI cho terminal
GREEN  = "\033[92m"
YELLOW = "\033[93m"
CYAN   = "\033[96m"
RED    = "\033[91m"
BOLD   = "\033[1m"
RESET  = "\033[0m"

def load_intent_model():
    print(f"  Đang tải Intent Classifier từ: {settings.INTENT_MODEL_REPO_ID}")
    df = pd.read_csv(settings.INTENT_DATASET_URL)
    le = LabelEncoder()
    le.fit(df["label"])

    tokenizer = AutoTokenizer.from_pretrained("vinai/phobert-base")
    try:
        model = AutoModelForSequenceClassification.from_pretrained(
            settings.INTENT_MODEL_REPO_ID
        )
    except Exception:
        model = AutoModelForSequenceClassification.from_pretrained(
            settings.INTENT_MODEL_REPO_ID, subfolder="intent_model_output"
        )
    model.eval()
    print("  Intent Classifier loaded.\n")
    return tokenizer, model, le

def load_ner_model():
    print(f"  Đang tải NER Model từ: {settings.NER_MODEL_REPO_ID}")
    ner = hf_pipeline(
        "token-classification",
        model=settings.NER_MODEL_REPO_ID,
        tokenizer=settings.NER_MODEL_REPO_ID,
        aggregation_strategy="simple",
    )
    print("  NER Model loaded.\n")
    return ner

def predict_intent(text: str, tokenizer, model, label_encoder):
    segmented = word_tokenize(text, format="text")
    inputs = tokenizer(
        segmented,
        return_tensors="pt",
        truncation=True,
        padding=True,
        max_length=128,
    )
    with torch.no_grad():
        logits = model(**inputs).logits

    pred_id    = torch.argmax(logits).item()
    intent     = label_encoder.inverse_transform([pred_id])[0]
    confidence = torch.softmax(logits, dim=1)[0][pred_id].item()
    return intent, confidence, segmented



def predict_ner(text: str, ner_pipe):
    raw = ner_pipe(text)
    if not raw:
        return []

    merged = []
    for ent in raw:
        group = ent.get("entity_group", ent.get("entity", "UNKNOWN"))
        word  = ent["word"].replace("@@", "").strip()
        score = ent["score"]

        if merged and merged[-1]["group"] == group:
            prev_raw = merged[-1]["_raw_word"]
            if ent["word"].startswith("@@") or prev_raw.endswith("@@"):
                merged[-1]["word"] += word
            else:
                merged[-1]["word"] += " " + word
            merged[-1]["_raw_word"] = ent["word"]
        else:
            merged.append({"group": group, "word": word, "_raw_word": ent["word"], "score": score})

    return [{"group": e["group"], "word": e["word"], "score": e["score"]} for e in merged]



def print_result(user_input: str, intent: str, confidence: float,
                 segmented: str, entities: list):
    route  = INTENT_ROUTE_MAP.get(intent, f"  {intent.upper()}")
    bar    = "─" * 62

    print(f"\n{bar}")
    print(f"  {BOLD}Câu hỏi :{RESET} {user_input}")
    print(f"  {BOLD}Tách từ :{RESET} {CYAN}{segmented}{RESET}")
    print(bar)

    # Confidence color
    conf_pct = confidence * 100
    color = GREEN if conf_pct >= 80 else (YELLOW if conf_pct >= 60 else RED)
    print(f"  {BOLD}ROUTE   :{RESET} {BOLD}{route}{RESET}")
    print(f"  {BOLD}Tự tin  :{RESET} {color}{conf_pct:.1f}%{RESET}")

    if entities:
        print(f"  {BOLD}Entities:{RESET}")
        for ent in entities:
            ent_color = CYAN
            print(f"       {ent_color}{ent['group'].ljust(28)}{RESET}: \"{BOLD}{ent['word']}{RESET}\"  ({ent['score']*100:.0f}%)")
    else:
        print(f"  {BOLD}Entities:{RESET} (không có)")

    print(bar)



def main():
    print("\n" + "=" * 62)
    print("  CHATBOT PIPELINE TEST  —  Intent + NER")
    print("=" * 62)
    print("  Đang khởi tạo các model, vui lòng chờ...\n")

    tokenizer, intent_model, label_encoder = load_intent_model()
    ner_pipe = load_ner_model()

    print("=" * 62)
    print("  Sẵn sàng! Nhập câu hỏi (gõ 'q' để thoát)")
    print("=" * 62)

    while True:
        try:
            user_input = input("\n👤 Khách hàng: ").strip()
        except (KeyboardInterrupt, EOFError):
            print("\nThoát.")
            break

        if user_input.lower() in ("q", "quit", "exit"):
            print("Thoát.")
            break

        if not user_input:
            continue

        intent, confidence, segmented = predict_intent(
            user_input, tokenizer, intent_model, label_encoder
        )
        entities = predict_ner(user_input, ner_pipe)
        print_result(user_input, intent, confidence, segmented, entities)


if __name__ == "__main__":
    main()
