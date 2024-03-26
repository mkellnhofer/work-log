function init() {
    document.querySelectorAll('.wl-delete-form').forEach(function(form) {
        form.onsubmit = function() {
            return confirm(form.dataset.dialogMessage);
        };
    });

    document.querySelectorAll('.wl-toggle-form').forEach(function(form) {
        input = form.getElementsByTagName('input')[1];
        input.onchange = function() {
            form.submit();
        };
    });
}

init();