// Supondo que seu backend Go responde em /api/users com [{id, nome}]
async function carregarUsuarios() {
    const res = await fetch('/api/users');
    const usuarios = await res.json();
    const lista = document.getElementById('user-list');
    lista.innerHTML = '';
    usuarios.forEach(usuario => {
        const card = document.createElement('div');
        card.className = 'user-card';
        card.innerHTML = `
            <div class="user-info">
                <span class="user-name">${usuario.nome}</span>
                <span class="user-id">ID: ${usuario.id}</span>
            </div>
            <div class="user-actions">
                <button class="user-action" onclick="enviarDM('${usuario.id}', '${usuario.nome}')">Enviar DM</button>
                <button class="user-action ban" onclick="banirUsuario('${usuario.id}', '${usuario.nome}')">Banir</button>
            </div>
        `;
        lista.appendChild(card);
    });
}

async function enviarDM(id, nome) {
    const mensagem = prompt(`Mensagem para ${nome}:`);
    if (!mensagem) return;
    alert('Função de envio de DM ainda não implementada no backend.');
}

async function banirUsuario(id, nome) {
    if (!confirm(`Tem certeza que deseja banir ${nome}?`)) return;
    alert('Função de banir ainda não implementada no backend.');
}

carregarUsuarios();