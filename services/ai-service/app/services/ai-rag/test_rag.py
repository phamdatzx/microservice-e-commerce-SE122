
import os
import sys

# Thêm đường dẫn ai-service gốc để thư mục 'app' được nhận dạng
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))))

from app.services.router import route_intent
from app.services.rag_service import ask_review_summary
from app.services.rag_service import ask_product_detail

print("="*50)
print("TEST 1: HỎI CHI TIẾT SẢN PHẨM")
question_1 = "Mô tả sản phẩm 'Camisa Denim Top Retro Mujer Estilo Hong Kong Tie Tencel Principios Otoño Slim-Fit Nuevo Pequeña' cho tôi"
print("Q:", question_1)
print("A:", ask_product_detail(question_1))

print("\n" + "="*50)
print("TEST 2: TÓM TẮT ĐÁNH GIÁ")
question_2 = "Nhận xét của khách hàng về sản phẩm 'Camisa Denim Top Retro Mujer Estilo Hong Kong Tie Tencel Principios Otoño Slim-Fit Nuevo Pequeña' như thế nào?"
print("Q:", question_2)
print("A:", ask_review_summary(question_2))
print("="*50)
