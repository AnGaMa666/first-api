document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("album-form");
    form.addEventListener("submit", function (event) {
        event.preventDefault();
        const formData = new FormData(form);
        fetch("http://localhost:8081/albums", {
            method: "POST",
            body: JSON.stringify(Object.fromEntries(formData)),
            headers: { "Content-Type": "application/json" }
        })
        .then(response => response.json())
        .then(data => {
            document.getElementById("response").innerText = "New Album Created: " + JSON.stringify(data);
            form.reset();
        })
        .catch(error => console.error(error));
    });
});
