document.addEventListener('DOMContentLoaded', function () {
    const addUserForm = document.getElementById('addUserForm');
    const userTableBody = document.querySelector('#users tbody');
    const resetUsersButton = document.getElementById('resetUsers');
    const startAttendanceButton = document.getElementById('startAttendance');

    // 获取用户列表
    function fetchUsers() {
        fetch('/api/users')
            .then(response => response.json())
            .then(data => {
                userTableBody.innerHTML = ''; // 清空表格
                data.forEach(user => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${user.id}</td>
                        <td>${user.qq}</td>
                        <td>${user.nickname}</td>
                        <td>${user.qq_name}</td>
                    `;
                    userTableBody.appendChild(row);
                });
            })
            .catch(error => console.error('Error fetching users:', error));
    }

    // 添加用户
    addUserForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const formData = new FormData(addUserForm);
        const userData = {
            id: formData.get('id'),
            qq: formData.get('qq'),
            level: parseInt(formData.get('nickname')),
            qq_name: formData.get('qq_name')
        };

        fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
                addUserForm.reset();
                fetchUsers(); // 重新获取用户列表
            })
            .catch(error => console.error('Error adding user:', error));
    });

  

    // 重置用户
    resetUsersButton.addEventListener('click', function () {
        fetch('/api/reset', {
            method: 'POST'
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
                fetchUsers(); // 重新获取用户列表
            })
            .catch(error => console.error('Error resetting users:', error));
    });

    // 开始考勤
    startAttendanceButton.addEventListener('click', function () {
        fetch('/api/attendance', {
            method: 'POST'
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
            })
            .catch(error => console.error('Error starting attendance:', error));
    });

    // 初始获取用户列表
    fetchUsers();
});
