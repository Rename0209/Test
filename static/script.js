function fetchData() {
    fetch('/data')
        .then(response => response.json())
        .then(data => {
            let tableBody = document.getElementById("dataTable");
            tableBody.innerHTML = "";
            data.forEach(item => {
                let row = `<tr>
                    <td>${item.topic_title}</td>
                    <td>${item.recipient_id}</td>
                    <td>${item.notification_messages_token}</td>
                    <td>${item.notification_messages_reoptin}</td>
                    <td>${item.user_token_status}</td>
                    <td>${new Date(item.creation_timestamp).toLocaleString()}</td>
                    <td>${new Date(item.token_expiry_timestamp).toLocaleString()}</td>
                    <td>${new Date(item.next_eligible_time * 1000).toLocaleString()}</td>
                </tr>`;
                tableBody.innerHTML += row;
            });
        })
        .catch(error => console.error('Lỗi tải dữ liệu:', error));
}
