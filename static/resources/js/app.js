const clearModal = () => {
  const modalContainer = document.getElementById('wl-modal-container');
  modalContainer.innerHTML = '';
}

// Clear modal when HTMX restores history (user presses back/forward)
document.body.addEventListener('htmx:historyRestore', clearModal);