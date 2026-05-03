import torch
import pandas as pd
from transformers import AutoTokenizer, AutoModelForSequenceClassification
from sklearn.preprocessing import LabelEncoder
import logging

from app.core.config import get_settings

logger = logging.getLogger(__name__)

settings = get_settings()

# Hugging Face Repository IDs and URLs
MODEL_REPO_ID = settings.INTENT_MODEL_REPO_ID
DATASET_URL = settings.INTENT_DATASET_URL

# ==============================
# 1. Load dataset để lấy label
# ==============================
logger.info(f"Loading dataset for label encoder from {DATASET_URL}")
try:
    df = pd.read_csv(DATASET_URL)
    label_encoder = LabelEncoder()
    label_encoder.fit(df["label"])
except Exception as e:
    logger.error(f"Failed to load dataset: {e}")
    raise e

# ==============================
# 2. Load model đã train
# ==============================
logger.info(f"Loading pretrained model from Hugging Face Hub: {MODEL_REPO_ID}")
try:
    tokenizer = AutoTokenizer.from_pretrained("vinai/phobert-base")
    model = AutoModelForSequenceClassification.from_pretrained(MODEL_REPO_ID)
    model.eval()
except Exception as e:
    logger.error(f"Failed to load model: {e}")
    raise e


# ==============================
# 3. Predict function
# ==============================
def predict_intent(text: str) -> str:
    """Predict the intent of a given text string."""
    inputs = tokenizer(
        text,
        return_tensors="pt",
        truncation=True,
        padding=True,
        max_length=128
    )

    with torch.no_grad():
        outputs = model(**inputs)

    logits = outputs.logits
    predicted_class_id = torch.argmax(logits).item()
    predicted_intent = label_encoder.inverse_transform([predicted_class_id])[0]

    return predicted_intent