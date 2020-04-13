import fetch from "isomorphic-unfetch";

export default async (req, res) => {
  let list = await fetch("http://storage:8080/list").then(res => res.json());
  res.end(JSON.stringify(list));
};
