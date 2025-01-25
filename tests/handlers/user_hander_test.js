// JS调用接口测试验证
// 放到浏览器控制台运行

// 测试：用户创建接口
fetch("http://127.0.0.1:8080/api/user", {
  method: "POST", // POST对应新建
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    "name": "用户1" + new Date().getTime(),
    "email": "user" + new Date().getTime() + "@xx.com",
    "phone": "124537890",
    "wechat": "mywechat",
    "address": "用户1 Address"
  }),
  mode: "cors"
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User created result:", data);
  })
  .catch(error => {
    console.error("Error creating user:", error);
  });

// 测试：更新用户
fetch("http://127.0.0.1:8080/api/user", {
  method: "PUT", // PUT对应更新
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    "id": 2,
    "name": "测试2",
    "email": "test2@xx.com",
    "phone": "1234537890",
    "wechat": "updated_wechat",
    "address": "Updted Address"
  }),
  mode: "cors"
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User updated result:", data);
  })
  .catch(error => {
    console.error("Error updating user:", error);
  });

// 测试：查询指定用户
fetch("http://127.0.0.1:8080/api/user/1", {
  method: "GET", // GET对应查询
  headers: {
    "Content-Type": "application/json",
  }
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("Users get result:", data);
  })
  .catch(error => {
    console.error("Error getting user:", error);
  });

// 测试：查询全部用户
fetch("http://127.0.0.1:8080/api/users", {
  method: "GET", // GET对应查询
  headers: {
    "Content-Type": "application/json",
  }
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User get result:", data);
  })
  .catch(error => {
    console.error("Error getting user:", error);
  });

// 测试：删除指定用户
fetch("http://127.0.0.1:8080/api/user/31", {
  method: "DELETE", // DELETE对应新建
  headers: {
    "Content-Type": "application/json",
  },
  mode: "cors"
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User updated result:", data);
  })
  .catch(error => {
    console.error("Error updating user:", error);
  });

  // 测试：获取分页
fetch(`http://127.0.0.1:8080/api/user/list?page=1&size=9&condition=email%20like%20"a%25"`, {
  method: "GET",
  headers: {
    "Content-Type": "application/json",
  },
  mode: "cors"
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User updated result:", data);
  })
  .catch(error => {
    console.error("Error updating user:", error);
  });

// 测试：获取分组
fetch(`http://127.0.0.1:8080/api/user/group?field=name,email`, {
  method: "GET",
  headers: {
    "Content-Type": "application/json",
  },
  mode: "cors"
})
.then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log("User updated result:", data);
  })
  .catch(error => {
    console.error("Error updating user:", error);
  });
  
/* 接口验证
#!/bin/bash
# 测试接口：用户列表
curl -X GET "http://127.0.0.1:8080/api/user/list?page=1&size=9&condition=email%20like%20%22a%25%22" \
     -H "Accept: application/json"
*/