const state = {
  seats: [],
  selectedSeatNo: "",
  myHold: null,
  timerId: null,
};

const els = {
  userId: document.getElementById("user-id"),
  refreshBtn: document.getElementById("refresh-btn"),
  seatGrid: document.getElementById("seat-grid"),
  holdBtn: document.getElementById("hold-btn"),
  payBtn: document.getElementById("pay-btn"),
  selectedSeatLabel: document.getElementById("selected-seat-label"),
  selectedSeatMeta: document.getElementById("selected-seat-meta"),
  detailsUser: document.getElementById("details-user"),
  detailsSeatID: document.getElementById("details-seat-id"),
  detailsToken: document.getElementById("details-token"),
  countdown: document.getElementById("countdown"),
  eventLog: document.getElementById("event-log"),
  connectionStatus: document.getElementById("connection-status"),
  seatTemplate: document.getElementById("seat-template"),
};

start();

async function start() {
  setupUser();
  wireEvents();
  await loadSeats();
  connectLiveFeed();
}

function setupUser() {
  const savedUser = localStorage.getItem("ticket-demo-user-id");
  els.userId.value = savedUser || `user-${Math.floor(Math.random() * 900 + 100)}`;
  saveUser();
}

function wireEvents() {
  els.userId.addEventListener("change", saveUser);
  els.refreshBtn.addEventListener("click", loadSeats);
  els.holdBtn.addEventListener("click", holdSeat);
  els.payBtn.addEventListener("click", payForSeat);
}

function saveUser() {
  localStorage.setItem("ticket-demo-user-id", els.userId.value.trim());
}

async function requestGraphQL(query) {
  const response = await fetch("/graphql", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query }),
  });

  const body = await response.json();
  if (body.errors?.length) {
    throw new Error(body.errors[0].message);
  }

  return body.data;
}

async function loadSeats() {
  try {
    const data = await requestGraphQL(`
      query {
        seats {
          id
          roomId
          number
          status
          heldBy
          holdToken
          expiresAt
        }
      }
    `);

    state.seats = data.seats;
    syncMyHold();
    renderSeats();
    renderSelection();
  } catch (error) {
    addFeedItem("Error", error.message);
  }
}

function renderSeats() {
  els.seatGrid.innerHTML = "";

  for (const seat of state.seats) {
    const card = els.seatTemplate.content.firstElementChild.cloneNode(true);
    card.querySelector(".seat-number").textContent = seat.number;
    card.querySelector(".seat-status").textContent = seat.status;
    card.classList.add(seat.status.toLowerCase());

    if (seat.number === state.selectedSeatNo) {
      card.classList.add("selected");
    }

    card.disabled = seat.status === "SOLD";
    card.addEventListener("click", () => {
      state.selectedSeatNo = seat.number;
      renderSeats();
      renderSelection();
    });

    els.seatGrid.appendChild(card);
  }
}

function renderSelection() {
  const seat = state.seats.find((item) => item.number === state.selectedSeatNo) || null;
  const myHold = state.myHold;

  els.selectedSeatLabel.textContent = seat ? `Seat ${seat.number}` : "Pick a seat";
  els.selectedSeatMeta.textContent = seat
    ? `Status: ${seat.status}${seat.heldBy ? ` | held by ${seat.heldBy}` : ""}`
    : "Choose any available seat to create a 5-minute lock.";

  els.detailsUser.textContent = myHold?.userId || "-";
  els.detailsSeatID.textContent = myHold?.seatId || seat?.id || "-";
  els.detailsToken.textContent = myHold?.holdToken || seat?.holdToken || "-";

  els.holdBtn.disabled = !(seat && seat.status === "AVAILABLE");
  els.payBtn.disabled = !(myHold && myHold.seatNumber === seat?.number);

  if (!myHold) {
    stopCountdown();
    els.countdown.textContent = "--:--";
  }
}

async function holdSeat() {
  const seat = state.seats.find((item) => item.number === state.selectedSeatNo);
  if (!seat) {
    return;
  }

  const userId = els.userId.value.trim() || "user-demo";

  try {
    const data = await requestGraphQL(`
      mutation {
        holdSeat(seatNumber: "${seat.number}", userId: "${userId}", amountCents: 12000) {
          id
          number
          status
          heldBy
          holdToken
          expiresAt
        }
      }
    `);

    state.myHold = {
      seatId: data.holdSeat.id,
      seatNumber: data.holdSeat.number,
      holdToken: data.holdSeat.holdToken,
      expiresAt: data.holdSeat.expiresAt,
      userId,
    };
    state.selectedSeatNo = data.holdSeat.number;

    startCountdown(data.holdSeat.expiresAt);
    addFeedItem("Seat held", `${data.holdSeat.number} locked for ${userId}`);
    await loadSeats();
  } catch (error) {
    addFeedItem("Hold failed", error.message);
    alert(error.message);
  }
}

async function payForSeat() {
  if (!state.myHold) {
    return;
  }

  try {
    const data = await requestGraphQL(`
      mutation {
        payForSeat(
          seatId: "${state.myHold.seatId}",
          holdToken: "${state.myHold.holdToken}",
          userId: "${state.myHold.userId}",
          amountCents: 12000
        ) {
          accepted
          message
          paymentRef
        }
      }
    `);

    addFeedItem("Payment queued", data.payForSeat.paymentRef);
    setTimeout(loadSeats, 1200);
  } catch (error) {
    addFeedItem("Payment failed", error.message);
    alert(error.message);
  }
}

function startCountdown(expiresAt) {
  stopCountdown();

  const endAt = new Date(expiresAt).getTime();
  if (!Number.isFinite(endAt)) {
    els.countdown.textContent = "--:--";
    return;
  }

  updateCountdown(endAt);
  state.timerId = setInterval(() => updateCountdown(endAt), 1000);
}

function updateCountdown(endAt) {
  const leftMs = endAt - Date.now();
  if (leftMs <= 0) {
    stopCountdown();
    els.countdown.textContent = "00:00";
    state.myHold = null;
    renderSelection();
    loadSeats();
    return;
  }

  const totalSeconds = Math.floor(leftMs / 1000);
  const minutes = String(Math.floor(totalSeconds / 60)).padStart(2, "0");
  const seconds = String(totalSeconds % 60).padStart(2, "0");
  els.countdown.textContent = `${minutes}:${seconds}`;
}

function stopCountdown() {
  if (state.timerId) {
    clearInterval(state.timerId);
    state.timerId = null;
  }
}

function syncMyHold() {
  if (!state.myHold) {
    return;
  }

  const currentSeat = state.seats.find((seat) => seat.id === state.myHold.seatId);
  if (!currentSeat) {
    state.myHold = null;
    stopCountdown();
    return;
  }

  if (currentSeat.status === "SOLD") {
    addFeedItem("Seat sold", `${currentSeat.number} payment completed`);
    state.myHold = null;
    stopCountdown();
    return;
  }

  if (currentSeat.status === "AVAILABLE") {
    addFeedItem("Hold expired", `${currentSeat.number} is available again`);
    state.myHold = null;
    stopCountdown();
    return;
  }

  state.myHold.expiresAt = currentSeat.expiresAt;
  startCountdown(currentSeat.expiresAt);
}

function connectLiveFeed() {
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const socket = new WebSocket(`${protocol}://${window.location.host}/graphql/subscribe`);

  socket.addEventListener("open", () => {
    els.connectionStatus.textContent = "Live updates connected";
    socket.send(JSON.stringify({ type: "connection_init" }));
  });

  socket.addEventListener("message", async (event) => {
    const body = JSON.parse(event.data);

    if (body.type === "connection_ack") {
      socket.send(JSON.stringify({
        id: "seat-feed",
        type: "subscribe",
        payload: {
          query: `subscription { seatUpdates(roomId:"world-tour-main-stage") { id number status heldBy holdToken expiresAt } }`,
        },
      }));
      return;
    }

    if (body.type === "next") {
      const seatUpdate = body.payload?.data?.seatUpdates;
      if (!seatUpdate) {
        return;
      }

      addFeedItem("Realtime update", `${seatUpdate.number} -> ${seatUpdate.status}`);
      await loadSeats();
    }
  });

  socket.addEventListener("close", () => {
    els.connectionStatus.textContent = "Live updates disconnected, retrying...";
    setTimeout(connectLiveFeed, 1500);
  });

  socket.addEventListener("error", () => {
    els.connectionStatus.textContent = "Live updates error";
  });
}

function addFeedItem(title, message) {
  const item = document.createElement("li");
  item.innerHTML = `<strong>${escapeHTML(title)}</strong><span>${escapeHTML(message)}</span>`;
  els.eventLog.prepend(item);

  while (els.eventLog.children.length > 8) {
    els.eventLog.removeChild(els.eventLog.lastChild);
  }
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}
