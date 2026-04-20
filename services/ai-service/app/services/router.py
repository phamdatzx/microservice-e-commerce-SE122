from app.services.predict_intent_service import predict_intent as get_intent
from app.services.rag_service import ask_product_detail, ask_review_summary

def route_intent(query: str):
    intent = get_intent(query)
    
    print(f"Detected intent: {intent} for query: {query}")

    # Taạm thời để vậy do route_intent chưa hoàn thiện
    if intent == "product_detail":
        return ask_product_detail(query)

    elif intent == "review_summary":
        return ask_review_summary(query)

    # Các intent khác bạn có thể implement sau
    elif intent == "similar_product":
        return f"Logic cho similar_product chưa được implement. Query: {query}"

    elif intent == "recommend_product":
        return f"Logic cho recommend_product chưa được implement. Query: {query}"

    elif intent == "compare_product":
        return f"Logic cho compare_product chưa được implement. Query: {query}"

    elif intent == "add_to_cart":
        return f"Logic cho add_to_cart chưa được implement. Query: {query}"

    else:
        # Fallback / General Question
        return f"Xin lỗi, tôi không chắc mình hiểu ý bạn. Vui lòng nói rõ hơn về sản phẩm (Intent trả về: {intent})"