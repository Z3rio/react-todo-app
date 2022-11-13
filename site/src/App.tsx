import "./App.css";

import { useEffect, useState } from "react";
import axios from "axios";

import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

function App() {
  const [todos, setTodos] = useState([]);

  useEffect(() => {
    axios
      .get("http://localhost:8080/api/getTodos")
      .then(async (resp) => {
        setTodos(await resp.data.todos);
      })
      .catch((err) => {
        if (err) {
          console.log(err);
        }
      });
  }, []);

  return (
    <div className="TodoApp">
      <h1>Todo App</h1>

      {JSON.stringify(todos)}

      <div className="todo-list">
        <div className="new-todo-item">
          <TextField type="text" label="Text" variant="standard" size="small" />
          <Button variant="outlined">Add todo</Button>
        </div>
      </div>
    </div>
  );
}

export default App;
