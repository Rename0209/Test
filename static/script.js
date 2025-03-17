function fetchData() {
    fetch('/data')
        .then(response => response.json())
        .then(data => {
            let tableBody = document.getElementById("dataTable");
            tableBody.innerHTML = "";

            data.forEach(item => {

                let row = `<tr>
                    <td>${item.topic_title || "N/A"}</td>
                    <td>${item.recipient_id || "N/A"}</td>
                    <td>${item.notification_messages_token || "N/A"}</td>
                    <td>${item.notification_messages_reoptin || "N/A"}</td>
                    <td>${item.user_token_status || "N/A"}</td>
                    <td>${item.creation_timestamp ? new Date(item.creation_timestamp).toLocaleString() : "Invalid Date"}</td>
                    <td>${item.token_expiry_timestamp ? new Date(item.token_expiry_timestamp).toLocaleString() : "Invalid Date"}</td>
                    <td>${item.next_eligible_time ? new Date(item.next_eligible_time * 1000).toLocaleString() : "Invalid Date"}</td>
                </tr>`;
                tableBody.innerHTML += row;
            });
        })
        .catch(error => console.error('Lỗi tải dữ liệu:', error));
}

