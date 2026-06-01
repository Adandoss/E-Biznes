import os
import random
import requests
import ollama
from fastapi import FastAPI
from pydantic import BaseModel

#3.0 #####################
# Inicjalizacja FastAPI
app = FastAPI()

OLLAMA_API_URL = os.getenv("OLLAMA_API_URL", "http://host.docker.internal:11434")
MODEL_NAME = os.getenv("MODEL_NAME", "llama3.1:8b")

#4.0 #####################
# 5 otwarć oraz zamknięć

OPENING_PHRASES = [
    "Hej! Witaj w naszym sklepie z ubraniami. W czym mogę pomóc?",
    "Cześć! Szukasz ubrań czy masz jakieś pytanie o sklep?",
    "Dzień dobry! Jestem asystentem sklepowym. Co Cię dzisiaj interesuje?",
    "Witaj! Chętnie pomogę Ci znaleźć jakieś fajne ciuchy.",
    "Cześć! Masz jakieś pytania odnośnie ubrań lub zakupów?"
]

CLOSING_PHRASES = [
    "Dzięki za rozmowę! Miłego dnia i udanych zakupów.",
    "Mam nadzieję, że pomogłem. Do zobaczenia!",
    "Dzięki za kontakt. W razie czego jestem tutaj.",
    "Do usłyszenia! Udanych zakupów życzę.",
    "Super, polecam się na przyszłość. Trzymaj się!"
]

GREETING_KEYWORDS = ["cześć", "hej", "witaj", "dzien dobry", "dzień dobry", "siema", "witam"]
FAREWELL_KEYWORDS = ["pa", "do widzenia", "narazie", "na razie", "koniec", "dzięki", "dziękuje", "dziękuję", "bye"]

# 4.5 #####################
# foltrujemy zaadnienia
def is_shop_related(text: str) -> bool:
    text_lower = text.lower()
    forbidden_words = ["programowanie", "python", "matematyka", "oblicz", "fizyka", "historia", "gotowanie", "przepis"]
    for word in forbidden_words:
        if word in text_lower:
            return False
    return True

#5.0 #####################
# sentyment 
def analyze_sentiment(text: str) -> str:
    text_lower = text.lower()
    positive_words = ["super", "świetny", "polecam", "ładny", "dobry", "miły", "pomocny"]
    negative_words = ["zły", "brzydki", "słaby", "drogo", "błąd", "kiepski", "problem"]
    
    score = 0
    for w in positive_words:
        score += text_lower.count(w)
    for w in negative_words:
        score -= text_lower.count(w)
        
    return "positive" if score >= 0 else "negative"

def fetch_database_products() -> list:
    try:
        response = requests.get("http://backend:8080/products", timeout=5)
        if response.status_code == 200:
            return response.json()
    except Exception as e:
        print(f"Failed to fetch products from backend: {e}", flush=True)
    return []

class ChatRequest(BaseModel):
    message: str

class ChatResponse(BaseModel):
    response: str

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/chat", response_model=ChatResponse)
def chat(req: ChatRequest):
    user_msg = req.message.strip()
    msg_lower = user_msg.lower()

    has_greeting = any(word in msg_lower for word in GREETING_KEYWORDS)
    has_farewell = any(word in msg_lower for word in FAREWELL_KEYWORDS)

    #4.5 #####################
    # Blokada wejściowa dla zapytań spoza zagadnień sklepowych
    if not is_shop_related(user_msg):
        return ChatResponse(response="Przepraszam, ale pomagam tylko w tematach związanych z modą, ubraniami i zakupami w naszym sklepie.")

    products = fetch_database_products()
    products_context = ""
    
    if products:
        products_context = "Nasze rzeczywiste produkty w bazie danych:\n"
        for p in products:
            name = p.get("name", "Produkt")
            desc = p.get("description", "")
            price = p.get("price", 0.0)
            category = p.get("category", {}).get("name", "Odzież")
            products_context += f"- {name} (Kategoria: {category}) - Opis: {desc} - Cena: {price} PLN\n"
    else:
        products_context = "Obecnie brak produktów w bazie danych sklepu."

    #4.5 #####################
    # Prompt systemowy dla ollamy wymuszający ograniczenie tematów wyłącznie do ubrań i sklepu
    system_prompt = (
        "Jesteś asystentem sklepu odzieżowego. Odpowiadaj TYLKO na pytania o nasz sklep, "
        "ubrania, modę, rozmiary, ceny i zakupy. "
        "Oto lista RZECZYWISTYCH produktów w naszym sklepie pobrana z bazy danych:\n"
        f"{products_context}\n"
        "Gdy użytkownik zapyta o ceny, jakie mamy produkty lub który produkt jest najtańszy/najdroższy, "
        "MUSISZ odpowiedzieć dokładnie na podstawie powyższej listy rzeczywistych produktów. "
        "Wskaż poprawnie najtańszy produkt na bazie podanych cen. "
        "Jeśli pytanie dotyczy czegoś innego niż moda i nasz sklep, grzecznie odmów. Odpowiadaj krótko i rzeczowo po polsku."
    )

    try:
        #3.0 #####################
        # Inicjalizacja połączenia i przesłanie zapytania do llamy

        client = ollama.Client(host=OLLAMA_API_URL)
        res = client.generate(
            model=MODEL_NAME,
            prompt=user_msg,
            system=system_prompt,
            options={"temperature": 0.7}
        )
        reply = res.get("response", "").strip()
    except Exception as e:
        print(f"ERROR: Ollama call failed (Model: {MODEL_NAME}, Host: {OLLAMA_API_URL}). Exception: {e}", flush=True)
        if products:
            cheapest = min(products, key=lambda x: x.get("price", 999999.0))
            names_list = ", ".join([p.get("name") for p in products])
            reply = f"Mamy w ofercie: {names_list}. Najtańszy produkt to {cheapest.get('name')} w cenie {cheapest.get('price')} PLN!"
        else:
            reply = "Sklep odzieżowy zaprasza! Mamy super ubrania w dobrych cenach."

    #5.0 #####################
    # Weryfikacja sentymentu odpowiedzi 
    
    user_sentiment = analyze_sentiment(user_msg)
    reply_sentiment = analyze_sentiment(reply)
    if user_sentiment == "negative" or reply_sentiment == "negative":
        reply = "Przepraszam za kłopot. Chętnie pomogę Ci z uśmiechem w wyborze ubrań i odpowiem na pytania!"

    #4.0 #####################
    # Prependowanie losowego otwarcia lub appendowanie zamknięcia rozmowy

    if has_greeting:
        opening = random.choice(OPENING_PHRASES)
        reply = f"{opening}\n\n{reply}"
        
    if has_farewell:
        closing = random.choice(CLOSING_PHRASES)
        reply = f"{reply}\n\n{closing}"

    return ChatResponse(response=reply)
