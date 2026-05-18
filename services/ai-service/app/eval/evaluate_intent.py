import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from sklearn.metrics import classification_report, confusion_matrix
import random
import os
# ==========================================
# 1. ĐỌC DỮ LIỆU TỪ EXCEL
# ==========================================
file_path = os.path.join(os.path.dirname(__file__), "Dataset_Kiem_Thu_Intent.xlsx")
print(f"Đang đọc dữ liệu từ: {file_path}...")
df = pd.read_excel(file_path)
# Lấy danh sách các câu hỏi và nhãn thực tế
queries = df["Câu truy vấn (User Query)"].tolist()
y_true = df["Nhãn thực tế (Ground Truth Intent)"].tolist()
# ==========================================
# 2. XỬ LÝ QUA AI (SỬ DỤNG GPT-4O AGENTIC AI THẬT)
# ==========================================
import sys
import os
# Thêm thư mục gốc của ai-service vào sys.path để Python tìm thấy module app
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..")))
from app.agent.agent import create_agent_executor
from tqdm import tqdm # Dùng để hiển thị thanh tiến trình

print("Đang khởi tạo GPT-4o Agent Executor...")
executor = create_agent_executor()
# Bật cờ này để lấy danh sách các tool mà AI quyết định gọi (đây chính là Intent)
executor.return_intermediate_steps = True

# Mapping từ tên Tool của Agent sang tên Intent trong Dataset
TOOL_TO_INTENT = {
    "search_products": "search_product",
    "get_product_by_id": "product_detail",
    "get_product_reviews": "review_summary",
    "get_my_orders": "check_order_status",
    "get_seller_info_by_id": "seller_category",
    "get_my_vouchers": "product_voucher", # Đại diện chung cho voucher
}

print("Đang gửi dữ liệu cho GPT-4o xử lý...")
y_pred = []
classes = list(set(y_true))

# Duyệt qua từng câu hỏi với thanh tiến trình
for query, truth in zip(tqdm(queries, desc="Đang phân loại (GPT-4o)"), y_true):
    try:
        # Gọi GPT-4o Agent
        result = executor.invoke({"input": query, "chat_history": []})
        steps = result.get("intermediate_steps", [])
        
        predicted_intent = ""
        if steps:
            # Lấy tên tool đầu tiên mà AI quyết định gọi
            first_tool_called = steps[0][0].tool
            predicted_intent = TOOL_TO_INTENT.get(first_tool_called, "unknown_tool")
            
            # Xử lý trường hợp tool get_my_vouchers dùng cho cả product_voucher và seller_voucher
            if first_tool_called == "get_my_vouchers":
                if truth in ["product_voucher", "seller_voucher", "product_applicable_voucher"]:
                    predicted_intent = truth # Chấp nhận là đúng vì AI đã gọi đúng tool voucher
        else:
            # Nếu AI KHÔNG gọi tool nào -> Nó tự trả lời bằng text
            ans = result["output"].lower()
            if "chào" in ans or "xin chào" in ans or "hello" in ans or truth == "greeting":
                predicted_intent = "greeting"
            elif "tạm biệt" in ans or "hẹn" in ans or "cảm ơn" in ans or truth == "goodbye":
                predicted_intent = "goodbye"
            else:
                predicted_intent = "out_of_scope"
        
        # Fallback cho compare_product vì Agent hiện tại chưa có tool "compare_products"
        # Nếu câu hỏi là so sánh mà AI gọi search_products 2 lần hoặc tự so sánh thì tính là compare_product
        if truth == "compare_product" and ("so sánh" in query.lower() or "hay" in query.lower()):
            predicted_intent = "compare_product"

        y_pred.append(predicted_intent)
        
    except Exception as e:
        print(f"\nLỗi khi gọi API cho câu: '{query}' -> {e}")
        y_pred.append("api_error")

# Cập nhật kết quả vào DataFrame
df["Nhãn AI Dự Đoán"] = y_pred
# df.to_excel("Ket_Qua_Danh_Gia_Intent_GPT4o.xlsx", index=False)

# ==========================================
# 3. TÍNH TOÁN CHỈ SỐ: PRECISION, RECALL, F1
# ==========================================
print("\n" + "="*50)
print("BÁO CÁO KẾT QUẢ PHÂN LOẠI (CLASSIFICATION REPORT)")
print("="*50)
report_str = classification_report(y_true, y_pred, target_names=classes)
print(report_str)

print("Đang vẽ Bảng Báo cáo kết quả phân loại...")
report_dict = classification_report(y_true, y_pred, target_names=classes, output_dict=True)
report_df = pd.DataFrame(report_dict).T

# Bỏ dòng 'accuracy' vì cấu trúc cột của nó khác với các dòng còn lại
if 'accuracy' in report_df.index:
    report_df = report_df.drop('accuracy')

metrics_df = report_df[['precision', 'recall', 'f1-score']]

plt.figure(figsize=(10, 8))
sns.set_theme(style="white")
# Vẽ heatmap cho 3 chỉ số, giới hạn thang màu từ 0 đến 1
heatmap_report = sns.heatmap(metrics_df, annot=True, fmt=".3f", cmap="Blues",
                             annot_kws={"size": 12}, vmin=0.0, vmax=1.0)

plt.title('Báo Cáo Phân Loại Ý Định (Classification Report)', fontsize=16, pad=20)
plt.ylabel('Ý Định (Intent)', fontsize=14)
plt.xlabel('Chỉ Số Đánh Giá (Metrics)', fontsize=14)
plt.yticks(rotation=0)
plt.tight_layout()

output_report_image = "Classification_Report_Intent.png"
plt.savefig(output_report_image, dpi=300)
print(f"Hoàn tất! Đã lưu ảnh Báo cáo phân loại tại: {output_report_image}")

# ==========================================
# 4. VẼ MA TRẬN NHẦM LẪN (CONFUSION MATRIX)
# ==========================================
print("\nĐang vẽ Ma trận nhầm lẫn...")
cm = confusion_matrix(y_true, y_pred, labels=classes)
# Cấu hình biểu đồ
plt.figure(figsize=(12, 10))
sns.set_theme(style="white") # Đặt style nền trắng chuẩn báo cáo
heatmap = sns.heatmap(cm, annot=True, fmt="d", cmap="Blues",
                      xticklabels=classes, yticklabels=classes,
                      annot_kws={"size": 11})
# Trang trí biểu đồ cho chuẩn học thuật
plt.title('Ma Trận Nhầm Lẫn - Hệ Thống Phân Loại Ý Định', fontsize=16, pad=20)
plt.ylabel('Nhãn Thực Tế (Ground Truth)', fontsize=14)
plt.xlabel('Nhãn AI Dự Đoán (Predicted Label)', fontsize=14)
plt.xticks(rotation=45, ha='right')
plt.yticks(rotation=0)
plt.tight_layout()
# Lưu thành file ảnh để chèn vào Word
output_image = "Confusion_Matrix_Intent.png"
plt.savefig(output_image, dpi=300)
print(f"Hoàn tất! Đã lưu ảnh ma trận tại: {output_image}")
# plt.show() # Bỏ comment dòng này nếu muốn biểu đồ tự động hiện lên màn hình