import { BaseAPI, route } from "../base";

export interface PantryProduct {
  id: string;
  name: string;
  quantity: number;
  minQuantity: number;
  unit: string;
  expiryDate: string;
  purchasePrice: number;
  onShoppingList: boolean;
}

export interface PantryOverview {
  products: PantryProduct[];
  lowStock: PantryProduct[];
  expiring: PantryProduct[];
  shoppingList: PantryProduct[];
}

export interface PantryProductCreate {
  name: string;
  quantity: number;
  minQuantity: number;
  unit: string;
  expiryDate: string;
  purchasePrice: number;
}

export interface PantryProductUpdate extends PantryProductCreate {
  onShoppingList: boolean;
}

export interface PantryConsume {
  amount: number;
}

export class PantryAPI extends BaseAPI {
  overview() {
    return this.http.get<PantryOverview>({ url: route("/pantry") });
  }

  createProduct(body: PantryProductCreate) {
    return this.http.post<PantryProductCreate, PantryProduct>({
      url: route("/pantry/products"),
      body,
    });
  }

  updateProduct(id: string, body: PantryProductUpdate) {
    return this.http.put<PantryProductUpdate, PantryProduct>({
      url: route(`/pantry/products/${id}`),
      body,
    });
  }

  consume(id: string, amount = 1) {
    return this.http.post<PantryConsume, PantryProduct>({
      url: route(`/pantry/products/${id}/consume`),
      body: { amount },
    });
  }

  toggleShopping(id: string) {
    return this.http.post<void, PantryProduct>({
      url: route(`/pantry/products/${id}/shopping`),
    });
  }
}
