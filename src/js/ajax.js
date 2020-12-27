function ajax() {
  $.getJSON('https://localhost:8080', (data) => {
    console.log(data);
  })
}
