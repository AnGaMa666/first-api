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
            // Aktualisiere die Albumliste nach dem Hinzufügen
            updateAlbumList();
            
        })
        .catch(error => console.error(error));
    });

    const showListButton = document.getElementById("show-list");
    showListButton.addEventListener("click", function (event) {
        event.preventDefault();
        // Aktualisiere die Albumliste, ohne die DELETE-Anfrage auszulösen
        updateAlbumList();
    });

    const deleteButton = document.getElementById("delete-button");
    deleteButton.addEventListener("click", function (event) {
        event.preventDefault();

        const checkboxes = document.querySelectorAll('input[type="checkbox"]:checked');

        if (checkboxes.length === 0) {
            // Keine Checkboxen ausgewählt, also nicht löschen
            alert("Please select at least one item to delete.");
            return;
        }

        const selectedAlbumIds = Array.from(checkboxes).map(checkbox => checkbox.getAttribute("data-id"));

        fetch("http://localhost:8081/albums", {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ ids: selectedAlbumIds })
        })
        .then(response => {
            if (response.status === 200) {
                console.log("Selected items deleted");
                // Aktualisiere die Albumliste nach dem Löschen
                updateAlbumList();
            } else if (response.status === 404) {
                console.error("One or more items not found");
            } else {
                console.error("Error deleting selected items");
            }
        })
        .catch(error => console.error(error));
    });

    // Funktion zum Aktualisieren der Albumliste
    function updateAlbumList() {
        fetch("http://localhost:8081/albums", {
            method: "GET",
            headers: { "Content-Type": "application/json" }
        })
        .then(response => response.json())
        .then(data => {
            const formattedData = data.map(item => ({
                title: item.title,
                artist: item.artist,
                price: parseFloat(item.price).toLocaleString("de-DE", { style: "currency", currency: "EUR" }),
                id: item.id
            }));

            const table = `
                <table>
                    <tr>
                        <th>Check</th>
                        <th>ID</th>
                        <th>Title</th>
                        <th>Artist</th>
                        <th>Price</th>
                    </tr>
                    ${formattedData.map(item => `
                        <tr>
                            <td><input type="checkbox" data-id="${item.id}"></td>
                            <td>${item.id}</td>
                            <td>${item.title}</td>
                            <td>${item.artist}</td>
                            <td>${item.price}</td>
                        </tr>
                    `).join('')}
                </table>
            `;

            const jsonResponse = document.getElementById("json-response");
            jsonResponse.innerHTML = table;
            jsonResponse.style.textAlign = "center";
        })
        .catch(error => console.error(error));
    }
});
