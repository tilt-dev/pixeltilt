import request from "request-promise-native";

export default async (req, res) => {
  const {
    query: { name }
  } = req;

  const options = {
    method: "POST",
    url: "http://storage/read",
    json: {
      Name: name
    }
  };

  var img;
  try {
    let result = await request(options);
    img = Buffer.from(result.Body, "base64");
  } catch (err) {
    console.log("Error reading from storage");
    console.log(err);
    res.statusCode = 500;
    res.end();
    return;
  }

  res.writeHead(200, {
    "Content-Type": "image/png",
    "Content-Length": img.length
  });
  res.end(img);
};
