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
    const res = await fetch('/api/send-dm', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({id, mensagem})
    });
    const data = await res.json();
    if (data.success) {
        alert('DM enviada!');
    } else {
        alert('Erro ao enviar DM.');
    }
}

async function banirUsuario(id, nome) {
    if (!confirm(`Tem certeza que deseja banir ${nome}?`)) return;
    const res = await fetch('/api/ban-user', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({id})
    });
    const data = await res.json();
    if (data.success) {
        alert('Usuário banido!');
        carregarUsuarios();
    } else {
        alert('Erro ao banir usuário.');
    }
}

carregarUsuarios();