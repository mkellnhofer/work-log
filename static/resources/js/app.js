const startDownload = (downloadUrl) => {
  window.location.href = downloadUrl;
}

// Start file download when downloadFile event is triggered
document.body.addEventListener('downloadFile', (evt) => startDownload(evt.detail.value));

const clearModal = () => {
  const modalContainer = document.getElementById('wl-modal-container');
  modalContainer.innerHTML = '';
}

// Clear modal when HTMX restores history (user presses back/forward)
document.body.addEventListener('htmx:historyRestore', clearModal);