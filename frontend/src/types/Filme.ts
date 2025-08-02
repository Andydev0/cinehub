// Define a estrutura de um filme, espelhando o JSON que o backend Go retorna.
export interface Filme {
    id: number;
    titulo: string;
    sinopse: string;
    dataLancamento: string;
    caminhoPoster: string;
    notaMedia: number;
  }