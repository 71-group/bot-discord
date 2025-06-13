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
    <button class="user-action" onclick="alert('Função em breve!')">Enviar DM</button>
    <button class="user-action" onclick="alert('Função em breve!')">Banir</button>
`;

    userList.appendChild(card);
});

const feedback = document.getElementById('feedback');
if (data.status == "ok") {
    feedback.textContent = 'Mensagem enviada com sucesso!';
    feedback.style.color = 'green';
    document.getElementById('mensagem').value = '';
} else {
    feedback.textContent = 'Erro ao enviar mensagem: ' + data.error;
    feedback.style.color = 'red';
}