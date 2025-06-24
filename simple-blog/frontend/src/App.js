import React, { useState, useEffect } from "react";
import axios from "axios";
import ReactMarkdown from "react-markdown";

export default function App() {
	const [token, setToken] = useState("");
	const [articles, setArticles] = useState([]);
	const [view, setView] = useState("list"); // list | detail | create
	const [currentArticle, setCurrentArticle] = useState(null);

	const [loginForm, setLoginForm] = useState({ username: "", password: "" });
	const [createForm, setCreateForm] = useState({ title: "", content: "" });

	useEffect(() => {
		fetchArticles();
	}, []);

	const fetchArticles = async () => {
		const res = await axios.get("http://localhost:8080/articles");
		setArticles(res.data);
	};

	const handleLogin = async () => {
		try {
			const res = await axios.post("http://localhost:8080/login", loginForm);
			setToken(res.data.token);
			alert("登录成功");
		} catch {
			// 打印错误信息
			console.error("登录失败" + JSON.stringify(loginForm));
			alert("登录失败");
		}
	};

	const handleCreateArticle = async () => {
		try {
			await axios.post(
				"http://localhost:8080/articles",
				createForm,
				{ headers: { Authorization: token } }
			);
			alert("创建成功");
			setView("list");
			fetchArticles();
		} catch {
			alert("创建失败");
		}
	};

	const showDetail = (article) => {
		setCurrentArticle(article);
		setView("detail");
	};

	return (
		<div style={{ maxWidth: 700, margin: "auto", padding: 20 }}>
			<h1>简易Markdown博客</h1>
			{!token && (
				<div>
					<h3>登录</h3>
					<input
						placeholder="用户名"
						value={loginForm.username}
						onChange={(e) =>
							setLoginForm({ ...loginForm, username: e.target.value })
						}
					/>
					<br />
					<input
						type="password"
						placeholder="密码"
						value={loginForm.password}
						onChange={(e) =>
							setLoginForm({ ...loginForm, password: e.target.value })
						}
					/>
					<br />
					<button onClick={handleLogin}>登录</button>
				</div>
			)}

			{token && view === "list" && (
				<>
					<button onClick={() => setView("create")}>写文章</button>
					<h3>文章列表</h3>
					<ul>
						{articles.map((a) => (
							<li key={a.ID}>
								<a href="#!" onClick={() => showDetail(a)}>
									{a.Title}
								</a>
							</li>
						))}
					</ul>
				</>
			)}

			{view === "detail" && currentArticle && (
				<div>
					<button onClick={() => setView("list")}>返回列表</button>
					<h2>{currentArticle.Title}</h2>
					<ReactMarkdown>{currentArticle.Content}</ReactMarkdown>
				</div>
			)}

			{view === "create" && (
				<div>
					<button onClick={() => setView("list")}>返回列表</button>
					<h3>新建文章</h3>
					<input
						placeholder="标题"
						value={createForm.title}
						onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
					/>
					<br />
					<textarea
						rows={10}
						placeholder="Markdown内容"
						value={createForm.content}
						onChange={(e) => setCreateForm({ ...createForm, content: e.target.value })}
						style={{ width: "100%" }}
					/>
					<br />
					<button onClick={handleCreateArticle}>提交</button>
				</div>
			)}
		</div>
	);
}
