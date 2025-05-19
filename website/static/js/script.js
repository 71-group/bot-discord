const usuarios = [
    { nome: "João Silva", email: "joao@example.com" },
    { nome: "Maria Oliveira", email: "maria@example.com" },
    { nome: "Carlos Souza", email: "carlos@example.com" },
    { nome: "Fernanda Lima", email: "fernanda@example.com" },
    { nome: "Ricardo Matos", email: "ricardo@example.com" }
];

const userList = document.getElementById("userList");

usuarios.forEach(usuario => {
    const card = document.createElement("div");
    card.className = "user-card";

    card.innerHTML = `
    <div class="user-name">${usuario.nome}</div>
    <div class="user-email">${usuario.email}</div>
    `;

    userList.appendChild(card);
});