document.addEventListener("DOMContentLoaded", function () {
    // Stellt sicher, dass der Code erst ausgeführt wird, wenn die DOM vollständig geladen ist.

    const form = document.getElementById("album-form");
    // Holt das HTML-Formular mit der ID "album-form".

    form.addEventListener("submit", function (event) {
        // Fügt einen "submit"-Event-Listener zum Formular hinzu.

        event.preventDefault();
        // Verhindert das Standardverhalten des Formulars, um die Seite nicht neu zu laden.

        const formData = new FormData(form);
        // Erstellt ein FormData-Objekt aus den Formulardaten.

        fetch("http://localhost:8081/albums", {
            // Sendet eine POST-Anfrage an den Server, um ein neues Album zu erstellen.
            method: "POST",
            body: JSON.stringify(Object.fromEntries(formData)),
            // Konvertiert die Formulardaten in JSON und sendet sie im Anfragekörper.
            headers: { "Content-Type": "application/json" }
            // Setzt den Header für JSON-Daten.
        })
        .then(response => response.json())
        // Parst die JSON-Antwort des Servers.

        .then(data => {
            // Verarbeitet die Antwort des Servers nach dem Erstellen des Albums.
            document.getElementById("response").innerText = "New Album Created: " + JSON.stringify(data);
            // Zeigt eine Erfolgsmeldung mit den Details des erstellten Albums an.
            form.reset();
            // Setzt das Formular zurück, um neue Eingaben zu ermöglichen.

            // Aktualisiere die Albumliste nach dem Hinzufügen
            updateAlbumList();
        })
        .catch(error => console.error(error));
        // Zeigt Fehlermeldungen in der Konsole an, wenn etwas schief geht.
    });

    const showListButton = document.getElementById("show-list");
    // Holt den HTML-Button mit der ID "show-list".

    showListButton.addEventListener("click", function (event) {
        // Fügt einen "click"-Event-Listener zum Button hinzu.

        event.preventDefault();
        // Verhindert das Standardverhalten des Buttons, um die Seite nicht neu zu laden.

        // Aktualisiere die Albumliste, ohne die DELETE-Anfrage auszulösen
        updateAlbumList();
        // Ruft die Funktion auf, um die Liste der Alben anzuzeigen.
    });

    const deleteButton = document.getElementById("delete-button");
    // Holt den HTML-Button mit der ID "delete-button".

    deleteButton.addEventListener("click", function (event) {
        // Fügt einen "click"-Event-Listener zum Button hinzu.

        event.preventDefault();
        // Verhindert das Standardverhalten des Buttons, um die Seite nicht neu zu laden.

        const checkboxes = document.querySelectorAll('input[type="checkbox"]:checked');
        // Sammelt alle ausgewählten Checkboxen, die zum Löschen markiert sind.

        if (checkboxes.length === 0) {
            // Überprüft, ob mindestens eine Checkbox ausgewählt ist.
            alert("Please select at least one item to delete.");
            // Zeigt eine Warnung, wenn keine Checkbox ausgewählt ist.
            return;
        }

        const selectedAlbumIds = Array.from(checkboxes).map(checkbox => checkbox.getAttribute("data-id"));
        // Extrahiert die IDs der ausgewählten Alben.

        fetch("http://localhost:8081/albums", {
            method: "DELETE",
            // Sendet eine DELETE-Anfrage an den Server, um ausgewählte Alben zu löschen.
            headers: { "Content-Type": "application/json" },
            // Setzt den Header für JSON-Daten.
            body: JSON.stringify({ ids: selectedAlbumIds })
            // Konvertiert die ausgewählten IDs in JSON und sendet sie im Anfragekörper.
        })
        .then(response => {
            if (response.status === 200) {
                console.log("Selected items deleted");
                // Zeigt eine Erfolgsmeldung in der Konsole an.
                // Aktualisiere die Albumliste nach dem Löschen
                updateAlbumList();
            } else if (response.status === 404) {
                console.error("One or more items not found");
                // Zeigt eine Fehlermeldung in der Konsole an, wenn Alben nicht gefunden wurden.
            } else {
                console.error("Error deleting selected items");
                // Zeigt eine allgemeine Fehlermeldung in der Konsole an.
            }
        })
        .catch(error => console.error(error));
        // Zeigt Fehlermeldungen in der Konsole an, wenn etwas schief geht.
    });

    // Funktion zum Aktualisieren der Albumliste
    function updateAlbumList() {
        fetch("http://localhost:8081/albums", {
            // Sendet eine GET-Anfrage an den Server, um die Liste der Alben abzurufen.
            method: "GET",
            headers: { "Content-Type": "application/json" }
            // Setzt den Header für JSON-Daten.
        })
        .then(response => response.json())
        // Parst die JSON-Antwort des Servers.

        .then(data => {
            // Verarbeitet die Antwort des Servers nach dem Aktualisieren der Liste.
            const formattedData = data.map(item => ({
                title: item.title,
                artist: item.artist,
                price: parseFloat(item.price).toLocaleString("de-DE", { style: "currency", currency: "EUR" }),
                id: item.id
            }));
            // Formatiert die Alben-Daten für die Anzeige.

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
            // Erstellt eine HTML-Tabelle mit den Alben-Daten.

            const jsonResponse = document.getElementById("json-response");
            jsonResponse.innerHTML = table;
            // Aktualisiert den Inhalt des HTML-Elements mit der ID "json-response".
            jsonResponse.style.textAlign = "center";
            // Zentriert den Text in diesem Element.
        })
        .catch(error => console.error(error));
        // Zeigt Fehlermeldungen in der Konsole an, wenn etwas schief geht.
    }
});
