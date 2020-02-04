import multer from "multer";
import util from "util";
import fs from "fs";
import request from "request-promise-native";

const uploadStorage = multer({ dest: "uploads/" });

export default async (req, res) => {
  let upload = util.promisify(uploadStorage.any());
  await upload(req, res);
  let filters = req.body.filters;
  let filename = req.files[0].filename;
  let filePath = "uploads/" + filename;

  let formData = {
    myFile: fs.createReadStream(filePath)
  };
  let parsedFilters = JSON.parse(filters);
  Object.keys(parsedFilters).forEach(item => {
    let key = item;
    let value = parsedFilters[item];

    formData[key] = value.toString();
  });

  const options = {
    method: "POST",
    url: "http://muxer/upload",
    headers: {
      "Content-Type": "multipart/form-data"
    },
    formData: formData,
    json: true
  };

  const proxyResult = await request(options);

  res.end(JSON.stringify(proxyResult));
};

export const config = {
  api: {
    bodyParser: false
  }
};
